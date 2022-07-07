package function

import (
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/{exch}/{sym:[A-Z0-9/-]+}", alertHandler)
	functions.HTTP("AlertHandler", r.ServeHTTP)

	validateEnv()
}

const (
	DefaultOpenOrderPercent      = "10"
	DefaultReducePositionPercent = "50"
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
	for _, ev := range []string{"FTX_APIKEY", "FTX_SECRET"} {
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
	var exch Exchange
	switch e := strings.ToUpper(v["exch"]); e {
	case "FTX":
		exch = NewFtx()
	case "BINANCE":
		fallthrough
	default:
		http.Error(w, "unsupported exchange:"+e, http.StatusBadRequest)
		return
	}

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	signal := NewAlert(string(bs))
	fmt.Println("action: ", signal.String(), sym)
	switch signal {
	case LONG:
		err = exch.LongPosition(sym)
	case REDUCE:
		err = exch.ReducePosition(sym)
	case CLOSE:
		err = exch.ClosePosition(sym)
	case STOP_LOSS:
		err = exch.StopLossPosition(sym)
	default:
		fmt.Println(signal, string(bs))
	}
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
