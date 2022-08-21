package data

import (
	"fmt"
	"strings"
)

const (
	UnknownAlert Alert = iota
	LONG
	SHORT
	REDUCE
	CLOSE
	STOP_LOSS
)

type Alert int

type AlertMessage struct {
	Strategy string
	Signal   string
	Kline    string
	Price    string
}

func NewAlertMessage(s string) (*AlertMessage, error) {
	s2 := strings.Split(s, "｜")
	if len(s2) != 4 {
		return nil, fmt.Errorf("unknown string format: %s", s)
	}

	return &AlertMessage{
		Strategy: s2[0],
		Signal:   s2[1],
		Kline:    s2[2],
		Price:    s2[3],
	}, nil
}

func NewAlert(s string) Alert {
	if len(s) == 0 {
		return UnknownAlert
	}

	msg, err := NewAlertMessage(s)
	if err != nil {
		fmt.Println(msg)
		return UnknownAlert
	}

	switch msg.Signal {
	case "空轉多訊號", "多方訊號", "空方平倉訊號":
		return LONG
	case "多轉空訊號", "空方訊號", "多方平倉訊號":
		return SHORT
	case "多方減倉訊號", "空方減倉訊號":
		return REDUCE
	case "多方停損訊號", "空方停損訊號":
		return STOP_LOSS
	default:
		fmt.Printf("unknown alert:%v len:%d\n", s, len(s))
		return UnknownAlert
	}
}

type Exchange interface {
	//MaxQuoteValue return order available to quote, e.g. free collateral * leverage
	//Spot instrument might return without leverage
	//sym provided for isolated margin account
	MaxQuoteValue(sym string) (total, free float64, err error)

	//GetPair return trading pair static info, which always can be cached
	GetPair(sym string) (*Pair, error)

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
	Type        MarketType
	Bid         float64
	Ask         float64
	Last        float64
	TickSize    float64
	MinNotional float64
}

type MarketType string

var (
	Spot   MarketType = "spot"
	Future MarketType = "future"
)

type Pair struct {
	Name  string
	Type  MarketType
	Base  string
	Quote string
}

func (p Pair) IsFuture() bool {
	return p.Type == Future
}

func (p Pair) IsSpot() bool {
	return p.Type == Spot
}
