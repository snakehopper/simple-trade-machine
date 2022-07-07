package function

import (
	"errors"
	"fmt"
	"ghohoo.solutions/yt/ftx"
	"ghohoo.solutions/yt/ftx/structs"
	"os"
	"strconv"
	"time"
)

type Client struct {
	*ftx.Client
}

func (c Client) LongPosition(sym string) error {
	acc, err := c.GetAccount()
	if err != nil {
		return err
	} else if !acc.Success {
		return errors.New("fetch account false success")
	}

	res := acc.Result

	pct, ok := os.LookupEnv("OPEN_PERCENT")
	if !ok {
		pct = DefaultOpenOrderPercent
	}
	pctValue, err := strconv.Atoi(pct)
	if err != nil {
		return err
	}

	orderUsd := (res.Collateral * res.Leverage) * (float64(pctValue) / 100)
	if freeUsd := res.FreeCollateral * res.Leverage * 0.9; freeUsd < orderUsd {
		fmt.Printf("low collateral, trim notional %v -> %v\n", orderUsd, freeUsd)
		orderUsd = freeUsd
	}

	//skip action if already has position
	for _, pos := range res.Positions {
		if pos.Future == sym {
			if pos.Size != 0 {
				fmt.Println("already has position, skip raise LONG again")
				return nil
			}
			break
		}
	}

	market, err := c.GetMarket(sym)
	if err != nil {
		return err
	} else if !market.Success {
		return errors.New("fetch market false success")
	}

	px := market.Result.Bid
	orderQty := orderUsd / px
	if orderQty < market.Result.MinProvideSize {
		fmt.Printf("order size too small (%v < %v), skip LONG action\n", orderQty, market.Result.MinProvideSize)
		return nil
	}

	resp, err := c.PlaceOrder(sym, "buy", px, "limit", orderQty, false, false, false)
	if err != nil {
		fmt.Printf("place order error: %v\n", err)
		return err
	} else if !resp.Success {
		return errors.New("place order unknown error")
	}

	po := resp.Result
	fmt.Printf("success place order. sym=%v px=%v qty=%v\n", po.Market, po.AvgFillPrice, po.Size)
	return nil
}

func (c Client) ReducePosition(sym string) error {
	pct, ok := os.LookupEnv("REDUCE_PERCENT")
	if !ok {
		pct = DefaultReducePositionPercent
	}
	pctValue, err := strconv.Atoi(pct)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		err = c.closePartialPosition(sym, pctValue)
		if err == nil {
			return nil
		}
		fmt.Printf("#%d wait a second and retry\n", i)
		time.Sleep(3 * time.Second)
	}
	return err
}

func (c Client) closePartialPosition(sym string, pct int) error {
	acc, err := c.GetAccount()
	if err != nil {
		return err
	} else if !acc.Success {
		return errors.New("fetch account false success")
	}

	var pos structs.Position
	for _, p := range acc.Result.Positions {
		if p.Future == sym {
			pos = p
			break
		}
	}
	fmt.Printf("%+v\n", pos)

	//skip action if EMPTY position
	if pos.Size == 0 {
		fmt.Println("empty position, skip close action")
		return nil
	}

	var offsetSide string
	if pos.Side == "buy" {
		offsetSide = "sell"
	} else if pos.Side == "sell" {
		offsetSide = "buy"
	}
	resp, err := c.PlaceOrder(sym, offsetSide, 0, "market", pos.Size*float64(pct)/100,
		true, true, false)
	if err != nil {
		fmt.Printf("close position error: %v\n", err)
		return err
	} else if !resp.Success {
		return errors.New("close position unknown error ")
	}

	po := resp.Result
	fmt.Printf("success close position. sym=%v px=%v size=%v\n", po.Market, po.AvgFillPrice, po.Size)
	return nil
}

func (c Client) ClosePosition(sym string) error {
	var err error
	for i := 0; i < 3; i++ {
		err = c.closePartialPosition(sym, 100)
		if err == nil {
			return nil
		}
		fmt.Printf("#%d wait a second and retry\n", i)
		time.Sleep(3 * time.Second)
	}
	return err
}

func (c Client) StopLossPosition(sym string) error {
	var err error
	for i := 0; i < 3; i++ {
		err = c.closePartialPosition(sym, 100)
		if err == nil {
			return nil
		}
		fmt.Printf("#%d wait a second and retry\n", i)
		time.Sleep(3 * time.Second)
	}
	return err
}

func NewFtx() *Client {
	return &Client{
		ftx.New(os.Getenv("FTX_APIKEY"), os.Getenv("FTX_SECRET"), os.Getenv("FTX_SUBACCOUNT")),
	}
}
