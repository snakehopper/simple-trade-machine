package spot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
)

func (a Api) IsolatedAccountInfo(sym ...string) (*MarginAccountResp, error) {
	var v = url.Values{}
	if len(sym) > 0 {
		v.Set("symbols", strings.Join(sym, ","))
	}

	resp, err := a.Get("/sapi/v1/margin/isolated/account", v, true)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out MarginAccountResp
	if err := json.Unmarshal(bs, &out); err != nil {
		return nil, err
	}

	if out.Code != 0 {
		return &out, fmt.Errorf("%v", string(bs))
	}

	return &out, nil
}
