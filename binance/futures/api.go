package futures

import (
	"encoding/json"
	"fmt"
	"ghohoo.solutions/yt/binance/com"
	"ghohoo.solutions/yt/internal/data"
	"io/ioutil"
	"net/url"
	"strings"
)

type Api struct {
	*com.Client
}

func NewApi(apiKey, secret string) *Api {
	return &Api{
		Client: com.NewClient("https://fapi.binance.com", apiKey, secret),
	}
}

func (a Api) MaxQuoteValue(sym string) (total, free float64, err error) {
	acc, err := a.AccountInfo()
	if err != nil {
		return
	} else if acc.Code != 0 {
		err = fmt.Errorf("fetch account error code:%v msg:%v", acc.Code, acc.Msg)
		return
	}

	sp, err := acc.GetPosition(sym)
	if err != nil {
		return
	}

	total = acc.TotalMarginBalance * sp.Leverage
	free = acc.AvailableBalance * sp.Leverage
	return
}

func (a Api) GetMarket(sym string) (*data.Market, error) {
	res, err := a.ExchangeInfo()
	if err != nil {
		return nil, err
	}

	for _, rs := range res.Symbols {
		if rs.Symbol == sym {
			var tickSize, minNotional float64
			for _, fp := range rs.Filters {
				if fp.FilterType == "MIN_NOTIONAL" {
					minNotional = fp.Notional
				} else if fp.FilterType == "PRICE_FILTER" {
					tickSize = fp.TickSize
				}
			}
			return &data.Market{
				Bid:         0, //TODO
				Ask:         0, //TODO
				Last:        0, //TODO
				TickSize:    tickSize,
				MinNotional: minNotional,
			}, nil
		}
	}
	return nil, fmt.Errorf("symbol %v not found", sym)
}

func (a Api) GetPosition(sym string) (float64, error) {
	acc, err := a.AccountInfo()
	if err != nil {
		return 0, err
	} else if acc.Code != 0 {
		err = fmt.Errorf("fetch account error code:%v msg:%v", acc.Code, acc.Msg)
		return 0, err
	}

	pos, err := acc.GetPosition(sym)
	if err != nil {
		return 0, err
	}
	return pos.PositionAmt, nil
}

func (a Api) LimitOrder(sym string, side data.Side, px float64, qty float64, ioc bool, _postOnly bool) error {
	exch, err := a.ExchangeInfo()
	if err != nil {
		return err
	}
	var v = url.Values{}
	v.Set("side", strings.ToUpper(string(side)))
	v.Set("symbol", sym)
	v.Set("type", "LIMIT")
	if ioc {
		v.Set("timeInForce", "IOC")
	} else {
		v.Set("timeInForce", "GTC")
	}
	rounded := exch.RoundLotSize(sym, qty)
	v.Set("quantity", fmt.Sprint(rounded))
	v.Set("price", fmt.Sprint(px))
	fmt.Println(v.Encode())
	resp, err := a.Post("/fapi/v1/order", v, true)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var out OrderResp
	if err := json.Unmarshal(bs, &out); err != nil {
		fmt.Println(err, string(bs))
		return err
	}

	if out.Code != 0 {
		return fmt.Errorf("limit order error: %v", out.Msg)
	}

	fmt.Println("<", string(bs))

	return nil
}

func (a Api) MarketOrder(sym string, side data.Side, quoteQty *float64, baseQty *float64) error {
	exch, err := a.ExchangeInfo()
	if err != nil {
		return err
	}

	var v = url.Values{}
	v.Set("side", strings.ToUpper(string(side)))
	v.Set("symbol", sym)
	v.Set("type", "MARKET")

	if baseQty != nil {
		rounded := exch.RoundLotSize(sym, *baseQty)
		v.Set("quantity", fmt.Sprint(rounded))
	} else {
		tk, err := a.OrderBookTicker(sym)
		if err != nil {
			return err
		}
		if side == data.Buy {
			rounded := exch.RoundLotSize(sym, *quoteQty/tk.AskPrice)
			v.Set("quantity", fmt.Sprint(rounded))
		} else {
			rounded := exch.RoundLotSize(sym, *quoteQty/tk.BidPrice)
			v.Set("quantity", fmt.Sprint(rounded))
		}
	}
	fmt.Println(v.Encode())
	resp, err := a.Post("/fapi/v1/order", v, true)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var out OrderResp
	if err := json.Unmarshal(bs, &out); err != nil {
		fmt.Println(err, string(bs))
		return err
	}

	if out.Code != 0 {
		return fmt.Errorf("market order error: %v", out.Msg)
	}

	fmt.Println("<", string(bs))

	return nil
}
