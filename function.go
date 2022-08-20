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
		if _, err := strconv.Atoi(v); err != nil {
			panic(err)
		}
	}
	if v, ok := os.LookupEnv("REDUCE_PERCENT"); ok {
		if _, err := strconv.Atoi(v); err != nil {
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
	var exch data.Exchange
	switch e {
	case "FTX":
		exch = ftx.NewApi(os.Getenv("FTX_APIKEY"), os.Getenv("FTX_SECRET"), os.Getenv("FTX_SUBACCOUNT"))
	case "BINANCE":
		exch = spot.NewApi(os.Getenv("BINANCE_APIKEY"), os.Getenv("BINANCE_SECRET"))
	case "BINANCE-FUTURES", "BINANCE_FUTURES", "FAPI":
		exch = futures.NewApi(os.Getenv("BINANCE_APIKEY"), os.Getenv("BINANCE_SECRET"))
	default:
		http.Error(w, "unsupported exchange:"+e, http.StatusBadRequest)
		return
	}

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	signal := data.NewAlert(string(bs))
	fmt.Println("action: ", signal.String(), sym)

	// lock to handle signal 1by1
	fmt.Println(filepath.Glob("/tmp/*"))
	fl := flock.New(fmt.Sprintf("/tmp/%s_%s", e, strings.ReplaceAll(sym, "/", "____")))
	if ok, err := fl.TryLock(); !ok || err != nil {
		fmt.Println("try lock failed", ok, err)
		fmt.Println("waiting lock with retry context...", e, sym)
		ok, err = fl.TryLockContext(context.Background(), time.Second)
		fmt.Println("TryLockContext return", ok, err)
	}
	defer fl.Unlock()

	switch signal {
	case data.LONG:
		err = longPosition(exch, sym)
	case data.SHORT:
		err = shortPosition(exch, sym)
	case data.REDUCE:
		err = reducePosition(exch, sym)
	case data.CLOSE:
		err = closePosition(exch, sym)
	case data.STOP_LOSS:
		err = stopLossPosition(exch, sym)
	default:
		fmt.Println(signal, string(bs))
	}
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func openPosition(exch data.Exchange, sym string, side data.Side) error {
	total, free, err := exch.MaxQuoteValue(sym)
	if err != nil {
		return nil
	}

	pct := getEnvFloat("OPEN_PERCENT", DefaultOpenOrderPercent)
	orderUsd := total * (pct / 100)
	if freeUsd := free * 0.9; freeUsd < orderUsd {
		fmt.Printf("low collateral, trim notional %v -> %v\n", orderUsd, freeUsd)
		orderUsd = freeUsd
	}

	market, err := exch.GetMarket(sym)
	if err != nil {
		return err
	}

	if orderUsd < market.MinNotional {
		fmt.Printf("order size too small (%v < %v), skip LONG action\n", orderUsd, market.MinNotional)
		return nil
	}

	return exch.MarketOrder(sym, side, &orderUsd, nil)
}

func closeIfAnyPosition(exch data.Exchange, sym string) error {
	pos, err := exch.GetPosition(sym)
	if err != nil {
		return err
	}
	if pos == 0 {
		return nil
	}

	return closePosition(exch, sym)
}

func longPosition(exch data.Exchange, sym string) error {
	if err := closeIfAnyPosition(exch, sym); err != nil {
		return err
	}
	return openPosition(exch, sym, data.Buy)
}

func shortPosition(exch data.Exchange, sym string) error {
	if err := closeIfAnyPosition(exch, sym); err != nil {
		return err
	}

	return openPosition(exch, sym, data.Sell)
}

func reducePosition(exch data.Exchange, sym string) error {
	pct := getEnvFloat("REDUCE_PERCENT", DefaultReducePositionPercent)

	var err error
	for i := 0; i < 3; i++ {
		err = closePartialPosition(exch, sym, pct)
		if err == nil {
			return nil
		}
		fmt.Printf("#%d wait a second and retry\n", i)
		time.Sleep(3 * time.Second)
	}
	return err
}

func closePartialPosition(exch data.Exchange, sym string, pct float64) error {
	pos, err := exch.GetPosition(sym)
	if err != nil {
		return err
	}
	//skip action if EMPTY position
	if pos == 0 {
		fmt.Println("empty position, skip close action")
		return nil
	}

	var offsetSide data.Side
	if pos > 0 {
		offsetSide = data.Sell
	} else {
		offsetSide = data.Buy
	}
	size := pos * pct / 100
	if err := exch.MarketOrder(sym, offsetSide, nil, &size); err != nil {
		fmt.Printf("close position error: %v\n", err)
		return err
	}

	return nil
}

func closePosition(exch data.Exchange, sym string) error {
	var err error
	for i := 0; i < 3; i++ {
		err = closePartialPosition(exch, sym, 100)
		if err == nil {
			return nil
		}
		fmt.Printf("#%d wait a second and retry\n", i)
		time.Sleep(3 * time.Second)
	}
	return err
}

func stopLossPosition(exch data.Exchange, sym string) error {
	var err error
	for i := 0; i < 3; i++ {
		err = closePartialPosition(exch, sym, 100)
		if err == nil {
			return nil
		}
		fmt.Printf("#%d wait a second and retry\n", i)
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
