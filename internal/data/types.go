package data

import "fmt"

const (
	UnknownAlert Alert = iota
	LONG
	REDUCE
	CLOSE
	STOP_LOSS
)

type Alert int

func NewAlert(s string) Alert {
	if len(s) == 0 {
		return UnknownAlert
	}

	switch a := []rune(s)[0]; a {
	case '多':
		return LONG
	case '減':
		return REDUCE
	case '平':
		return CLOSE
	case '停':
		return STOP_LOSS
	default:
		fmt.Printf("unknown alert:%v len:%d\n", a, len(s))
		return UnknownAlert
	}
}

type Exchange interface {
	//MaxQuoteValue return order available to quote, e.g. free collateral * leverage
	//sym provided for isolated margin account
	MaxQuoteValue(sym string) (total, free float64, err error) // e.g. collateral * leverage

	GetMarket(sym string) (*Market, error)

	GetPosition(sym string) (float64, error)
	LimitOrder(sym string, side Side, px float64, qty float64, ioc bool, postOnly bool) error
	MarketOrder(sym string, side Side, px *float64, qty *float64) error
}

type Side string

const (
	Buy  Side = "buy"
	Sell Side = "sell"
)

type Market struct {
	Bid         float64
	Ask         float64
	Last        float64
	TickSize    float64
	MinNotional float64
}
