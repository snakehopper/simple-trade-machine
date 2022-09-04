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
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
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
}

var whitelist = NewWhitelist()

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

	signal, err := NewSignal(string(bs))
	if err != nil {
		fmt.Printf("msg:[%s] err:%v\n", string(bs), err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log := setupLogger(e, sym, signal)

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

	h, err := NewSignalHandler(signal, e, sym, r.URL.Query())
	if err != nil {
		http.Error(w, "unsupported exchange:"+e, http.StatusBadRequest)
		return
	}

	switch signal.Action {
	case LONG:
		err = h.longPosition()
	case SHORT:
		err = h.shortPosition()
	case REDUCE:
		err = h.reducePosition()
	case CLOSE_LONG:
		err = h.closeIfAnyPositionNow(data.Buy)
	case CLOSE_SHORT:
		err = h.closeIfAnyPositionNow(data.Sell)
	case STOP_LOSS:
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
	log      *zap.SugaredLogger
	sig      *Signal
	override url.Values
	exch     data.Exchange
	sym      string
	exName   string
}

func NewSignalHandler(sig *Signal, exName, symbol string, q url.Values) (*SignalHandler, error) {
	logger := setupLogger(exName, symbol, sig)

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
		log:      logger,
		sig:      sig,
		override: q,
		sym:      symbol,
		exName:   exName,
		exch:     exch,
	}, nil
}

func (h SignalHandler) GetOrderType() data.OrderType {
	switch s := strings.ToLower(h.override.Get("order")); s {
	case "limit":
		return data.LimitOrder
	case "market":
		return data.MarketOrder
	}

	switch strings.ToLower(h.stringFromEnv("ORDER_TYPE")) {
	case "limit":
		return data.LimitOrder
	case "market":
		return data.MarketOrder
	default:
		panic("unreachable")
	}
}

func (h SignalHandler) stringFromEnv(k string) string {
	if v := h.override.Get(strings.ToLower(k)); v != "" {
		return v
	}

	k1 := fmt.Sprintf("%s_%s", h.sig.Strategy, k)

	if s := viper.GetString(k1); s != "" {
		return s
	}

	return viper.GetString(k)
}

func (h SignalHandler) floatFromEnv(k string) float64 {
	if s := h.override.Get(strings.ToLower(k)); s != "" {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			return v
		}
	}

	k1 := fmt.Sprintf("%s_%s", h.sig.Strategy, k)
	if viper.IsSet(k1) {
		return viper.GetFloat64(k1)
	}

	return viper.GetFloat64(k)
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

	pct := h.floatFromEnv("OPEN_PERCENT")
	orderUsd := total * (pct / 100)
	h.log.Infof("MaxQuoteValue(%s) total:%v free:%v orderUsd:%v", h.sym, total, free, orderUsd)

	if freeUsd := free * 0.9; freeUsd < orderUsd {
		h.log.Infof("low collateral, trim notional %v -> %v", orderUsd, freeUsd)
		orderUsd = freeUsd
	}

	if orderUsd < market.MinNotional {
		h.log.Infof("order size too small (%v < %v), skip %v action", orderUsd, market.MinNotional, side)
		return nil
	}

	if h.GetOrderType() == data.MarketOrder {
		_, err = h.exch.MarketOrder(h.sym, side, &orderUsd, nil)
		return err
	}

	px, err := h.bestQuotePrice(side)
	if err != nil {
		return err
	}

	oid, err := h.exch.LimitOrder(h.sym, side, px, orderUsd/px, false, false)
	if err != nil {
		return err
	}

	go h.followUpLimitOrder(oid)
	return nil
}

func (h SignalHandler) closeIfAnyPositionNow(holding data.Side) error {
	pos, err := h.exch.GetPosition(h.sym)
	if err != nil {
		return err
	}
	switch holding {
	case data.Buy:
		if pos > 0 {
			return h.closePosition(true)
		}
	case data.Sell:
		if pos < 0 {
			return h.closePosition(true)
		}
	}
	h.log.Info("no position to close")
	return nil
}

func (h SignalHandler) longPosition() error {
	err := h.closeIfAnyPositionNow(data.Sell)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		if err = h.openPosition(data.Buy); err == nil {
			return nil
		}
		h.log.Infof("#%d err:%v wait a second and retry", i, err)
		time.Sleep(3 * time.Second)
	}
	return err
}

func (h SignalHandler) shortPosition() error {
	if err := h.closeIfAnyPositionNow(data.Buy); err != nil {
		return err
	}

	if p, err := h.exch.GetPair(h.sym); err != nil {
		return err
	} else if p.IsSpot() {
		h.log.Info("Spot not to short")
		return nil
	}

	var err error
	for i := 0; i < 3; i++ {
		if err = h.openPosition(data.Sell); err == nil {
			return nil
		}
		h.log.Infof("#%d err:%v wait a second and retry", i, err)
		time.Sleep(3 * time.Second)
	}
	return err
}

func (h SignalHandler) reducePosition() error {
	pct := h.floatFromEnv("REDUCE_PERCENT")
	var err error
	for i := 0; i < 3; i++ {
		if err = h.closePartialPosition(pct, false); err == nil {
			return nil
		}
		h.log.Infof("#%d err:%v wait a second and retry", i, err)
		time.Sleep(3 * time.Second)
	}
	return err
}

func (h SignalHandler) closePartialPosition(pct float64, force bool) error {
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

	if force || h.GetOrderType() == data.MarketOrder {
		_, err = h.exch.MarketOrder(h.sym, offsetSide, nil, &size)
		return err
	}

	px, err := h.bestQuotePrice(offsetSide)
	if err != nil {
		return err
	}

	oid, err := h.exch.LimitOrder(h.sym, offsetSide, px, size, false, false)
	if err != nil {
		return err
	}

	go h.followUpLimitOrder(oid)
	return nil
}

func (h SignalHandler) closePosition(force bool) error {
	var err error
	for i := 0; i < 3; i++ {
		err = h.closePartialPosition(100, force)
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
		err = h.closePartialPosition(100, true)
		if err == nil {
			return nil
		}
		h.log.Infof("#%d err:%v wait a second and retry", i, err)
		time.Sleep(3 * time.Second)
	}
	return err
}

func (h SignalHandler) bestQuotePrice(side data.Side) (float64, error) {
	ob, err := h.exch.GetOrderBook(h.sym)
	if err != nil {
		return 0, err
	}

	var px float64
	if side == data.Buy {
		px = ob.Bid[0].Px
		if px*ob.Bid[0].Size < 100 {
			px = ob.Bid[1].Px
		}
	} else if side == data.Sell {
		px = ob.Ask[0].Px
		if px*ob.Ask[0].Size < 100 {
			px = ob.Ask[1].Px
		}
	}

	return px, nil
}

func (h SignalHandler) followUpLimitOrder(oid string) {
	d := viper.GetDuration("FOLLOWUP_LIMIT_ORDER")
	h.log.Infof("followup limit order %s after %v", oid, d)
	time.Sleep(d)

	od, err := h.exch.GetOrder(h.sym, oid)
	if err != nil {
		h.log.Warnf("followup limit order %s failed, err:%v", oid, err)
		return
	}

	if od.RemainingSize == 0 {
		h.log.Infof("limit order all filled")
		return
	}

	if err := h.exch.CancelOrder(h.sym, oid); err != nil {
		h.log.Warnf("cancel limit order %s failed, err:%v", oid, err)
		return
	}

	oid2, err := h.exch.MarketOrder(od.Pair.Name, od.Side, nil, &od.RemainingSize)
	h.log.Infof("limit-order:%s market-order:%s size:%v err:%v",
		oid, oid2, od.RemainingSize, err)
}
