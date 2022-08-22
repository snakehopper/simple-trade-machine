package ftx

import (
	"errors"
	"fmt"
	"ghohoo.solutions/yt/ftx/structs"
	"ghohoo.solutions/yt/internal/data"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"math"
	"strings"
)

type Api struct {
	log        *zap.SugaredLogger
	Api        string
	Secret     []byte
	Subaccount string
}

func (a Api) MaxQuoteValue(sym string) (total, free float64, err error) {
	acc, err := a.GetAccount()
	if err != nil {
		return
	} else if !acc.Success {
		err = errors.New("fetch account false success")
		return
	}

	res := acc.Result
	if a.GetTradingPair(sym).IsSpot() {
		mx := viper.GetFloat64("SPOT_OPEN_X")
		total = res.Collateral * mx
		free = res.FreeCollateral * mx
		return
	}

	total = res.Collateral * res.Leverage
	free = res.FreeCollateral * res.Leverage
	return
}

func (a Api) GetPosition(sym string) (float64, error) {
	switch a.GetTradingPair(sym).Type {
	case data.Spot:
		bal, err := a.GetBalance(sym)
		if err != nil {
			return 0, err
		}
		return bal.Free, nil
	case data.Future:
		acc, err := a.GetAccount()
		if err != nil {
			return 0, err
		} else if !acc.Success {
			return 0, errors.New("fetch account false success")
		}
		for _, p := range acc.Result.Positions {
			if p.Future == sym {
				return p.NetSize, nil
			}
		}
		fallthrough
	default:
		return 0, fmt.Errorf("unknown market: %s", sym)
	}
}

func (a Api) GetMarket(market string) (*data.Market, error) {
	resp, err := a._get("markets/"+market, []byte(""))
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
	var typ data.MarketType
	if res.Type == "spot" {
		typ = data.Spot
	} else if res.Type == "future" {
		typ = data.Future
	}

	return &data.Market{
		Bid:         res.Bid,
		Ask:         res.Ask,
		Last:        res.Last,
		Type:        typ,
		TickSize:    res.SizeIncrement,
		MinNotional: res.MinProvideSize,
	}, nil
}

func (a Api) GetPair(sym string) (*data.Pair, error) {
	p := a.GetTradingPair(sym)
	if p.Name == "" {
		return nil, fmt.Errorf("invalid symbol")
	}

	return &p, nil
}

func (a Api) LimitOrder(sym string, side data.Side, px float64, qty float64, ioc bool, postOnly bool) error {
	size := math.Abs(qty)
	resp, err := a.PlaceOrder(sym, strings.ToLower(string(side)), px, "limit", size, false, false, false)
	if err != nil {
		a.log.Infof("place limit order error: %v", err)
		return err
	} else if !resp.Success {
		return errors.New("place limit order unknown error")
	}

	return nil
}

func (a Api) MarketOrder(sym string, side data.Side, quoteUnit *float64, qty *float64) error {
	var size float64
	if qty != nil {
		size = math.Abs(*qty)
	} else if quoteUnit != nil {
		m, err := a.GetMarket(sym)
		if err != nil {
			return fmt.Errorf("get price error when MarketOrder: %w", err)
		}
		size = *quoteUnit / m.Last
	} else {
		return fmt.Errorf("either px or qty should defined")
	}
	resp, err := a.PlaceOrder(sym, strings.ToLower(string(side)), 0, "market", size,
		false, true, false)
	if err != nil {
		a.log.Infof("place market order error: %v", err)
		return err
	} else if !resp.Success {
		return errors.New("place order unknown error")
	}

	return nil
}

func NewApi(log *zap.SugaredLogger, apiKey, secret, subaccount string) *Api {
	return &Api{
		log:        log,
		Api:        apiKey,
		Secret:     []byte(secret),
		Subaccount: subaccount,
	}
}
