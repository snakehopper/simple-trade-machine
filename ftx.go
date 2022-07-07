package function

import (
	"errors"
	"fmt"
	"ghohoo.solutions/yt/ftx"
	"ghohoo.solutions/yt/ftx/structs"
	"log"
	"os"
	"strconv"
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
	if res.FreeCollateral < 50 {
		return nil
	}

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
	resp, err := c.PlaceOrder(sym, "buy", px, "limit", orderUsd/px, false, false, false)
	if err != nil {
		log.Printf("place order error: %v", err)
		return err
	} else if !resp.Success {
		return errors.New("place order false success")
	}

	po := resp.Result
	fmt.Printf("place order success. sym=%v px=%v qty=%v", po.Market, po.AvgFillPrice, po.Size)
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

	return c.closePartialPosition(sym, pctValue)
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
		true, false, false)
	if err != nil {
		log.Printf("close position error: %v", err)
		return err
	} else if !resp.Success {
		return errors.New("close position false success")
	}

	po := resp.Result
	fmt.Printf("close position success. sym=%v px=%v size=%v", po.Market, po.AvgFillPrice, po.Size)
	return nil
}

func (c Client) ClosePosition(sym string) error {
	return c.closePartialPosition(sym, 100)
}

func (c Client) StopLossPosition(sym string) error {
	return c.closePartialPosition(sym, 100)
}

func NewFtx() *Client {
	return &Client{
		ftx.New(os.Getenv("FTX_APIKEY"), os.Getenv("FTX_SECRET"), os.Getenv("FTX_SUBACCOUNT")),
	}
}
