package ftx

import (
	"fmt"
	"ghohoo.solutions/yt/ftx/structs"
	"ghohoo.solutions/yt/internal/data"
	"strconv"
)

type HistoricalPrices structs.HistoricalPrices
type Trades structs.Trades

func (a Api) GetHistoricalPrices(market string, resolution int64,
	limit int64, startTime int64, endTime int64) (HistoricalPrices, error) {
	var historicalPrices HistoricalPrices
	resp, err := a._get(
		"markets/"+market+
			"/candles?resolution="+strconv.FormatInt(resolution, 10)+
			"&limit="+strconv.FormatInt(limit, 10)+
			"&start_time="+strconv.FormatInt(startTime, 10)+
			"&end_time="+strconv.FormatInt(endTime, 10),
		[]byte(""))
	if err != nil {
		fmt.Printf("Error GetHistoricalPrices: %v\n", err)
		return historicalPrices, err
	}
	err = _processResponse(resp, &historicalPrices)
	return historicalPrices, err
}

func (a Api) GetTrades(market string, limit int64, startTime int64, endTime int64) (Trades, error) {
	var trades Trades
	resp, err := a._get(
		"markets/"+market+"/trades?"+
			"&limit="+strconv.FormatInt(limit, 10)+
			"&start_time="+strconv.FormatInt(startTime, 10)+
			"&end_time="+strconv.FormatInt(endTime, 10),
		[]byte(""))
	if err != nil {
		fmt.Printf("Error GetTrades: %v\n", err)
		return trades, err
	}
	err = _processResponse(resp, &trades)
	return trades, err
}

//GetTradingPair return cached symbol info
func (a Api) GetTradingPair(sym string) data.Pair {
	res, ok := khMarket[sym]
	if ok {
		return res
	}

	if err := a.FetchMarkets(); err != nil {
		a.log.Warnf("fetch markets error: %v", err)
		return data.Pair{}
	}

	return khMarket[sym]
}

func (a Api) FetchMarkets() error {
	resp, err := a._get("markets", []byte(""))
	if err != nil {
		return err
	}

	var markets = struct {
		Success bool             `json:"success"`
		Result  []structs.Market `json:"result"`
	}{}
	if err = _processResponse(resp, &markets); err != nil {
		return err
	}

	for _, m := range markets.Result {
		if m.Type == "spot" {
			khMarket[m.Name] = data.Pair{
				Name:  m.Name,
				Type:  data.Spot,
				Base:  *m.BaseCurrency,
				Quote: *m.QuoteCurrency,
			}
		} else if m.Type == "future" {
			khMarket[m.Name] = data.Pair{
				Name:  m.Name,
				Type:  data.Future,
				Base:  *m.Underlying,
				Quote: "USD",
			}
		}
	}

	return nil
}

var khMarket = make(map[string]data.Pair)
