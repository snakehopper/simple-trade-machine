package ftx

import (
	"fmt"
	"ghohoo.solutions/yt/ftx/structs"
	"ghohoo.solutions/yt/internal/data"
	"strconv"
)

type HistoricalPrices structs.HistoricalPrices
type Trades structs.Trades

func (client *Client) GetHistoricalPrices(market string, resolution int64,
	limit int64, startTime int64, endTime int64) (HistoricalPrices, error) {
	var historicalPrices HistoricalPrices
	resp, err := client._get(
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

func (client *Client) GetTrades(market string, limit int64, startTime int64, endTime int64) (Trades, error) {
	var trades Trades
	resp, err := client._get(
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

func (client *Client) GetMarket(market string) (*data.Market, error) {
	resp, err := client._get("markets/"+market, []byte(""))
	if err != nil {
		fmt.Printf("Error GetMarket: %v\n", err)
		return nil, err
	}
	var mResp structs.MarketResponse
	err = _processResponse(resp, &mResp)
	if err != nil {
		return nil, err
	}

	res := mResp.Result
	return &data.Market{
		Bid:         res.Bid,
		Ask:         res.Ask,
		Last:        res.Last,
		TickSize:    res.SizeIncrement,
		MinNotional: res.MinProvideSize,
	}, nil
}
