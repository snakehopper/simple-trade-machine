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

type MarketInfo struct {
	data map[string]Market
}

//GetTradingPair return cached symbol info
func (c *Client) GetTradingPair(sym string) Market {
	res, ok := c.markets[sym]
	if ok {
		return res
	}

	if err := c.FetchMarkets(); err != nil {
		fmt.Printf("fetch markets error: %v", err)
		return Market{}
	}

	return c.markets[sym]
}

func (c *Client) FetchMarkets() error {
	resp, err := c._get("markets", []byte(""))
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
			c.markets[m.Name] = Market{
				Name:  m.Name,
				Type:  Spot,
				Base:  *m.BaseCurrency,
				Quote: *m.QuoteCurrency,
			}
		} else if m.Type == "future" {
			c.markets[m.Name] = Market{
				Name:  m.Name,
				Type:  Future,
				Base:  *m.Underlying,
				Quote: "USD",
			}
		}
	}

	return nil
}

type Market struct {
	Name  string
	Type  MarketType
	Base  string
	Quote string
}

type MarketType string

var (
	Spot   MarketType = "spot"
	Future MarketType = "future"
)
