package spot

import (
	"encoding/json"
	"fmt"
	"ghohoo.solutions/yt/binance/com"
	"ghohoo.solutions/yt/internal/data"
	"go.uber.org/zap"
	"io/ioutil"
	"math"
	"net/url"
	"strings"
)

type Api struct {
	*com.Client
	log *zap.SugaredLogger
}

func NewApi(logger *zap.SugaredLogger, apiKey, secret string) *Api {
	return &Api{
		Client: com.NewClient(logger, "https://api.binance.com", apiKey, secret),
		log:    logger,
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

	quote := a.GetTradingPair(sym).Quote
	for _, bal := range acc.Balances {
		if bal.Asset == quote {
			free = bal.Free
			total = bal.Free + bal.Locked
			return
		}
	}
	return
}

func (a Api) GetMarket(sym string) (*data.Market, error) {
	res, err := a.ExchangeInfo(sym)
	if err != nil {
		return nil, err
	}

	for _, rs := range res.Symbols {
		if rs.Symbol == sym {
			var tickSize, minNotional float64
			for _, fp := range rs.Filters {
				if fp.FilterType == "MIN_NOTIONAL" {
					minNotional = fp.MinNotional
				} else if fp.FilterType == "PRICE_FILTER" {
					tickSize = fp.TickSize
				}
			}
			return &data.Market{
				Bid:         0, //TODO
				Ask:         0, //TODO
				Last:        0, //TODO
				Type:        data.Spot,
				TickSize:    tickSize,
				MinNotional: minNotional,
			}, nil
		}
	}
	return nil, fmt.Errorf("symbol %v not found", sym)
}

func (a Api) GetPair(sym string) (*data.Pair, error) {
	p := a.GetTradingPair(sym)
	if p.Name == "" {
		return nil, fmt.Errorf("invalid symbol")
	}

	return &p, nil
}

func (a Api) GetPosition(sym string) (float64, error) {
	acc, err := a.AccountInfo()
	if err != nil {
		return 0, err
	} else if acc.Code != 0 {
		err = fmt.Errorf("fetch account error code:%v msg:%v", acc.Code, acc.Msg)
		return 0, err
	}

	exch, err := a.ExchangeInfo(sym)
	if err != nil {
		return 0, err
	}
	var base string
	if s := exch.Symbols[0]; s.Symbol != sym {
		err = fmt.Errorf("exchange info symbol not match: %+v", s)
		return 0, err
	} else {
		base = s.BaseAsset
	}
	for _, bal := range acc.Balances {
		if bal.Asset == base {
			return bal.Free, nil
		}
	}

	return 0, fmt.Errorf("unknown symbol %v", sym)
}

func (a Api) LimitOrder(sym string, side data.Side, px float64, qty float64, ioc bool, _postOnly bool) error {
	exch, err := a.ExchangeInfo(sym)
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
	resp, err := a.Post("/api/v3/order", v, true)
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
		a.log.Infof("unmarshal OrderResp err:%v payload:%v", err, strings.TrimSpace(string(bs)))
		return err
	}

	if out.Code != 0 {
		return fmt.Errorf("market order error: %v", out.Msg)
	}

	a.log.Info("<", strings.TrimSpace(string(bs)))
	return nil
}

func (a Api) MarketOrder(sym string, side data.Side, quoteQty *float64, baseQty *float64) error {
	exch, err := a.ExchangeInfo(sym)
	if err != nil {
		return err
	}

	var v = url.Values{}
	v.Set("side", strings.ToUpper(string(side)))
	v.Set("symbol", sym)
	v.Set("type", "MARKET")
	if quoteQty != nil {
		rounded := exch.RoundTickSize(sym, *quoteQty)
		if rounded < exch.MinNotional(sym) {
			a.log.Info("rounded quoteOrderQty smaller than MIN_NOTIONAL, skip place order")
			return nil
		}
		v.Set("quoteOrderQty", fmt.Sprint(rounded))
	}
	if baseQty != nil {
		rounded := exch.RoundLotSize(sym, math.Abs(*baseQty))
		if rounded == 0 {
			a.log.Info("rounded quantity smaller than LOT_SIZE, skip place order")
			return nil
		}
		tk, err := a.OrderBookTicker(sym)
		if err != nil {
			return err
		}
		if rounded*tk.BidPrice < exch.MinNotional(sym) {
			a.log.Info("rounded quoteOrderQty smaller than MIN_NOTIONAL, skip place order")
			return nil
		}
		v.Set("quantity", fmt.Sprint(rounded))
	}
	resp, err := a.Post("/api/v3/order", v, true)
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
		return err
	}

	if out.Code != 0 {
		return fmt.Errorf("market order error: %v", out.Msg)
	}

	a.log.Info("<", strings.TrimSpace(string(bs)))
	return nil
}
