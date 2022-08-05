package function

import (
	"errors"
	"fmt"
	"ghohoo.solutions/yt/ftx"
	"ghohoo.solutions/yt/ftx/structs"
	"ghohoo.solutions/yt/internal/data"
	"math"
	"os"
)

type Client struct {
	*ftx.Client
}

func (c Client) MaxQuoteValue(sym string) (total, free float64, err error) {
	acc, err := c.GetAccount()
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

func (c Client) GetPosition(sym string) (float64, error) {
	acc, err := c.GetAccount()
	if err != nil {
		return 0, err
	} else if !acc.Success {
		return 0, errors.New("fetch account false success")
	}

	var pos structs.Position
	for _, p := range acc.Result.Positions {
		if p.Future == sym {
			pos = p
			break
		}
	}
	return pos.NetSize, nil
}

func (c Client) LimitOrder(sym string, side data.Side, px float64, qty float64, ioc bool, postOnly bool) error {
	size := math.Abs(qty)
	resp, err := c.PlaceOrder(sym, string(side), px, "limit", size, false, false, false)
	if err != nil {
		fmt.Printf("place limit order error: %v\n", err)
		return err
	} else if !resp.Success {
		return errors.New("place limit order unknown error")
	}

	return nil
}

func (c Client) MarketOrder(sym string, side data.Side, quoteUnit *float64, qty *float64) error {
	var size float64
	if qty != nil {
		size = math.Abs(*qty)
	} else if quoteUnit != nil {
		m, err := c.GetMarket(sym)
		if err != nil {
			return fmt.Errorf("get price error when MarketOrder: %w", err)
		}
		size = *quoteUnit / m.Last
	} else {
		return fmt.Errorf("either px or qty should defined")
	}
	resp, err := c.PlaceOrder(sym, string(side), 0, "market", size,
		false, true, false)
	if err != nil {
		fmt.Printf("place market order error: %v\n", err)
		return err
	} else if !resp.Success {
		return errors.New("place order unknown error")
	}

	return nil
}
func NewFtx() *Client {
	return &Client{
		ftx.New(os.Getenv("FTX_APIKEY"), os.Getenv("FTX_SECRET"), os.Getenv("FTX_SUBACCOUNT")),
	}
}
