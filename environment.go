package function

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

const (
	DefaultOpenOrderPercent      = 10.
	DefaultReducePositionPercent = 50.
	DefaultPositionOpenX         = 1.
)

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("OPEN_PERCENT", DefaultOpenOrderPercent)
	viper.SetDefault("REDUCE_PERCENT", DefaultReducePositionPercent)
	viper.SetDefault("SPOT_OPEN_X", DefaultPositionOpenX)

	// float part
	for _, ev := range []string{
		"OPEN_PERCENT", "REDUCE_PERCENT", "SPOT_OPEN_X",
	} {
		v := viper.GetString(ev)
		if _, err := strconv.ParseFloat(v, 64); err != nil {
			panic(err)
		}
	}

	//string part
	for _, ev := range []string{
		"FTX_APIKEY", "FTX_SECRET",
		"BINANCE_APIKEY", "BINANCE_SECRET",
	} {
		if _, ok := os.LookupEnv(ev); !ok {
			fmt.Println(ev, "not set")
		}
	}
}
