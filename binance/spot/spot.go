package spot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
)

func (a Api) ExchangeInfo(sym ...string) (*ExchangeInfoResp, error) {
	var v = url.Values{}
	if len(sym) == 1 {
		v.Set("symbol", sym[0])
	} else if len(sym) > 1 {
		bs, err := json.Marshal(sym)
		if err != nil {
			return nil, err
		}
		v.Set("symbols", string(bs))
	}

	resp, err := a.Get("/api/v3/exchangeInfo", v, false)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out ExchangeInfoResp
	if err := json.Unmarshal(bs, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (a Api) OrderBookTicker(sym string) (*OrderBookTickerResp, error) {
	q := url.Values{}
	q.Set("symbol", sym)
	resp, err := a.Get("/api/v3/ticker/bookTicker", q, false)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out OrderBookTickerResp
	if err := json.Unmarshal(bs, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (a Api) AccountInfo(sym ...string) (*AccountResp, error) {
	var v = url.Values{}
	if len(sym) > 0 {
		v.Set("symbols", strings.Join(sym, ","))
	}

	resp, err := a.Get("/api/v3/account", v, true)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out AccountResp
	if err := json.Unmarshal(bs, &out); err != nil {
		return nil, err
	}

	if out.Code != 0 {
		return &out, fmt.Errorf("%v", string(bs))
	}

	return &out, nil
}
