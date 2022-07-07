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
	r.HandleFunc("/{exch}/{sym:[A-Z0-9/-]+}", predAlert)
	functions.HTTP("PredAlert", r.ServeHTTP)

	validateEnv()
}

const (
	DefaultOpenOrderPercent      = "10"
	DefaultReducePositionPercent = "50"
)

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

func predAlert(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	sym := strings.Trim(v["sym"], "/")
	var exch Exchange
	switch e := strings.ToUpper(v["exch"]); e {
	case "FTX":
		exch = NewFtx()
	//	ok
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
	fmt.Println("action: ", signal.String())
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
		fmt.Println("unknown alert: ", signal)
	}
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
