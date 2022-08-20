package ftx

import (
	"errors"
	"fmt"
	"ghohoo.solutions/yt/internal/data"
	"math"
	"strings"
)

type Api struct {
	*Client
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
	total = res.Collateral * res.Leverage
	free = res.FreeCollateral * res.Leverage
	return
}

func (a Api) GetPosition(sym string) (float64, error) {
	switch a.GetTradingPair(sym).Type {
	case Spot:
		bal, err := a.GetBalance(sym)
		if err != nil {
			return 0, err
		}
		return bal.Free, nil
	case Future:
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

func (a Api) LimitOrder(sym string, side data.Side, px float64, qty float64, ioc bool, postOnly bool) error {
	size := math.Abs(qty)
	resp, err := a.PlaceOrder(sym, strings.ToLower(string(side)), px, "limit", size, false, false, false)
	if err != nil {
		fmt.Printf("place limit order error: %v\n", err)
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
		fmt.Printf("place market order error: %v\n", err)
		return err
	} else if !resp.Success {
		return errors.New("place order unknown error")
	}

	return nil
}
func NewApi(apiKey, secret, subaccount string) *Api {
	return &Api{
		New(apiKey, secret, subaccount),
	}
}
