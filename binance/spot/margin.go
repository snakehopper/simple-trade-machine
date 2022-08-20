package spot

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/url"
)

func (a Api) RepayMarginLoan(asset string, isolated bool, symbol string, amount float64) (*MarginRepayResp, error) {
	var v = url.Values{}
	v.Set("asset", asset)
	if isolated {
		v.Set("isIsolated", "TRUE")
		if symbol == "" {
			return nil, fmt.Errorf("symbol must be set when isolated")
		}
	}
	v.Set("symbol", symbol)
	v.Set("amount", decimal.NewFromFloat(amount).String())

	resp, err := a.Post("/sapi/v1/margin/repay", v, true)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out MarginRepayResp
	if err := json.Unmarshal(bs, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

type Side string
type SideEffect string

const (
	Buy          Side       = "BUY"
	Sell         Side       = "SELL"
	NoSideEffect SideEffect = "NO_SIDE_EFFECT"
	MarginBuy    SideEffect = "MARGIN_BUY"
	AutoRepay    SideEffect = "AUTO_REPAY"
)

func (a Api) MarginMarketOrder(isolated bool, sideEffect SideEffect, pair string, side Side, qty, quoteOrderQty *float64) (*OrderResp, error) {
	var v = url.Values{}
	v.Set("side", string(side))
	if isolated {
		v.Set("isIsolated", "TRUE")
		if pair == "" {
			return nil, fmt.Errorf("symbol must be set when isolated")
		}
	}
	v.Set("symbol", pair)
	v.Set("type", "MARKET")
	if qty != nil {
		v.Set("quantity", fmt.Sprint(*qty))
	}
	if quoteOrderQty != nil {
		v.Set("quoteOrderQty", fmt.Sprint(*quoteOrderQty))
	}
	v.Set("sideEffectType", string(sideEffect))
	resp, err := a.Post("/sapi/v1/margin/order", v, true)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out OrderResp
	if err := json.Unmarshal(bs, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (a Api) MarginAllOrders(isolated bool, sym string) ([]MarginOrderResp, error) {
	var v = url.Values{}
	if isolated {
		v.Set("isIsolated", "TRUE")
		if sym == "" {
			return nil, fmt.Errorf("symbol must be set when isolated")
		}
	}
	v.Set("symbol", sym)

	resp, err := a.Get("/sapi/v1/margin/allOrders", v, true)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out = make([]MarginOrderResp, 0)
	if err := json.Unmarshal(bs, &out); err != nil {
		return nil, err
	}

	return out, nil
}
