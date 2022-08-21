package function

import (
	"context"
	"fmt"
	"ghohoo.solutions/yt/binance/futures"
	"ghohoo.solutions/yt/binance/spot"
	"ghohoo.solutions/yt/ftx"
	"ghohoo.solutions/yt/internal/data"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gofrs/flock"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/{exch}/{sym:[A-Z0-9/-]+}", alertHandler)
	functions.HTTP("AlertHandler", r.ServeHTTP)

	validateEnv()
}

const (
	DefaultOpenOrderPercent      = 10.
	DefaultReducePositionPercent = 50.
)

var whitelist = NewWhitelist()

func validateEnv() {
	if v, ok := os.LookupEnv("OPEN_PERCENT"); ok {
		if _, err := strconv.ParseFloat(v, 64); err != nil {
			panic(err)
		}
	}
	if v, ok := os.LookupEnv("REDUCE_PERCENT"); ok {
		if _, err := strconv.ParseFloat(v, 64); err != nil {
			panic(err)
		}
	}
	for _, ev := range []string{
		"FTX_APIKEY", "FTX_SECRET",
		"BINANCE_APIKEY", "BINANCE_SECRET"} {
		if _, ok := os.LookupEnv(ev); !ok {
			fmt.Println(ev, "not set")
		}
	}
}

func alertHandler(w http.ResponseWriter, r *http.Request) {
	if !whitelist.Allow(r) {
		return
	}

	v := mux.Vars(r)
	sym := strings.Trim(v["sym"], "/")
	e := strings.ToUpper(v["exch"])

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	signal := data.NewAlert(string(bs))
	log := setupLogger(e, sym)
	log.Info("action: ", signal.String())

	// lock to handle signal 1by1
	log.Info(filepath.Glob("/tmp/*"))
	fl := flock.New(fmt.Sprintf("/tmp/%s_%s", e, strings.ReplaceAll(sym, "/", "____")))
	if ok, err := fl.TryLock(); !ok || err != nil {
		log.Infof("try lock failed. ok=%v err=%v", ok, err)
		log.Info("waiting lock with retry context...")
		ok, err = fl.TryLockContext(context.Background(), time.Second)
		log.Infof("TryLockContext return %v, %v", ok, err)
	}
	defer fl.Unlock()

	h, err := NewSignalHandler(e, sym)
	if err != nil {
		http.Error(w, "unsupported exchange:"+e, http.StatusBadRequest)
		return
	}

	switch signal {
	case data.LONG:
		err = h.longPosition()
	case data.SHORT:
		err = h.shortPosition()
	case data.REDUCE:
		err = h.reducePosition()
	case data.CLOSE:
		err = h.closePosition()
	case data.STOP_LOSS:
		err = h.stopLossPosition()
	default:
		log.Infof("unknown signal:%v data:'%s'", signal, string(bs))
	}
	if err != nil {
		log.Info(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

type SignalHandler struct {
	log    *zap.SugaredLogger
	exch   data.Exchange
	sym    string
	exName string
}

func NewSignalHandler(exName, symbol string) (*SignalHandler, error) {
	logger := setupLogger(exName, symbol)

	var exch data.Exchange
	switch exName {
	case "FTX":
		exch = ftx.NewApi(logger, os.Getenv("FTX_APIKEY"), os.Getenv("FTX_SECRET"), os.Getenv("FTX_SUBACCOUNT"))
	case "BINANCE":
		exch = spot.NewApi(logger, os.Getenv("BINANCE_APIKEY"), os.Getenv("BINANCE_SECRET"))
	case "BINANCE-FUTURES", "BINANCE_FUTURES", "FAPI":
		exch = futures.NewApi(logger, os.Getenv("BINANCE_APIKEY"), os.Getenv("BINANCE_SECRET"))
	default:
		return nil, fmt.Errorf("unsupported exchange")
	}

	return &SignalHandler{
		log:    logger,
		sym:    symbol,
		exName: exName,
		exch:   exch,
	}, nil
}

func (h SignalHandler) openPosition(side data.Side) error {
	market, err := h.exch.GetMarket(h.sym)
	if err != nil {
		return err
	}

	total, free, err := h.exch.MaxQuoteValue(h.sym)
	if err != nil {
		return nil
	}

	pct := getEnvFloat("OPEN_PERCENT", DefaultOpenOrderPercent)
	orderUsd := total * (pct / 100)
	if freeUsd := free * 0.9; freeUsd < orderUsd {
		h.log.Infof("low collateral, trim notional %v -> %v", orderUsd, freeUsd)
		orderUsd = freeUsd
	}

	if orderUsd < market.MinNotional {
		h.log.Infof("order size too small (%v < %v), skip %v action", orderUsd, side, market.MinNotional)
		return nil
	}

	return h.exch.MarketOrder(h.sym, side, &orderUsd, nil)
}

func (h SignalHandler) closeIfAnyPosition() error {
	pos, err := h.exch.GetPosition(h.sym)
	if err != nil {
		return err
	}
	if pos == 0 {
		h.log.Info("no position to close")
		return nil
	}

	return h.closePosition()
}

func (h SignalHandler) longPosition() error {
	if err := h.closeIfAnyPosition(); err != nil {
		return err
	}
	return h.openPosition(data.Buy)
}

func (h SignalHandler) shortPosition() error {
	if err := h.closeIfAnyPosition(); err != nil {
		return err
	}

	if p, err := h.exch.GetPair(h.sym); err != nil {
		return err
	} else if p.IsSpot() {
		h.log.Info("Spot not to short")
		return nil
	}

	return h.openPosition(data.Sell)
}

func (h SignalHandler) reducePosition() error {
	pct := getEnvFloat("REDUCE_PERCENT", DefaultReducePositionPercent)

	var err error
	for i := 0; i < 3; i++ {
		err = h.closePartialPosition(pct)
		if err == nil {
			return nil
		}
		h.log.Infof("#%d wait a second and retry", i)
		time.Sleep(3 * time.Second)
	}
	return err
}

func (h SignalHandler) closePartialPosition(pct float64) error {
	pos, err := h.exch.GetPosition(h.sym)
	if err != nil {
		return err
	}
	//skip action if EMPTY position
	if pos == 0 {
		h.log.Info("empty position, skip close action")
		return nil
	}

	var offsetSide data.Side
	if pos > 0 {
		offsetSide = data.Sell
	} else {
		offsetSide = data.Buy
	}
	size := pos * pct / 100
	if err := h.exch.MarketOrder(h.sym, offsetSide, nil, &size); err != nil {
		h.log.Infof("close position error: %v", err)
		return err
	}

	return nil
}

func (h SignalHandler) closePosition() error {
	var err error
	for i := 0; i < 3; i++ {
		err = h.closePartialPosition(100)
		if err == nil {
			return nil
		}
		h.log.Infof("#%d err:%v, wait a second and retry", i, err)
		time.Sleep(3 * time.Second)
	}
	return err
}

func (h SignalHandler) stopLossPosition() error {
	var err error
	for i := 0; i < 3; i++ {
		err = h.closePartialPosition(100)
		if err == nil {
			return nil
		}
		h.log.Infof("#%d err:%v wait a second and retry", i, err)
		time.Sleep(3 * time.Second)
	}
	return err
}

func getEnvFloat(key string, _default float64) float64 {
	res, ok := os.LookupEnv(key)
	if !ok {
		return _default
	}

	val, err := strconv.ParseFloat(res, 64)
	if err != nil {
		return _default
	}

	return val
}
