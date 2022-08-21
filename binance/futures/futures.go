package futures

import (
	"encoding/json"
	"fmt"
	"ghohoo.solutions/yt/internal/data"
	"io/ioutil"
	"net/url"
	"strings"
)

func (a Api) ExchangeInfo() (*ExchangeInfoResp, error) {
	resp, err := a.Get("/fapi/v1/exchangeInfo", url.Values{}, false)
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
	resp, err := a.Get("/fapi/v1/ticker/bookTicker", q, false)
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

	resp, err := a.Get("/fapi/v2/account", v, true)
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
	ex, err := a.ExchangeInfo()
	if err != nil {
		return err
	}

	for _, s := range ex.Symbols {
		khMarket[s.Symbol] = data.Pair{
			Type:  data.Future,
			Name:  s.Symbol,
			Base:  s.BaseAsset,
			Quote: s.QuoteAsset,
		}
	}

	return nil
}

var khMarket = make(map[string]data.Pair)
