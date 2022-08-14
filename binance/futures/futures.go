package futures

import (
	"encoding/json"
	"fmt"
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
