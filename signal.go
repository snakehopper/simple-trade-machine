//go:generate stringer -type Action -output signal_string.go
package function

import (
	"fmt"
	"strings"
	"time"
)

const (
	UnknownSignal Action = iota
	LONG
	SHORT
	REDUCE
	CLOSE_LONG
	CLOSE_SHORT
	STOP_LOSS

	CounterTrading  Strategy = "COUNTER"
	TrendFollowing  Strategy = "TREND"
	UnknownStrategy Strategy = "UnKNOwn"
)

type Action int
type Strategy string

type Signal struct {
	Strategy  Strategy
	Action    Action
	Triggered time.Time
}

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

func NewSignal(s string) (*Signal, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("empty alert message")
	}

	msg, err := NewAlertMessage(s)
	if err != nil {
		return nil, err
	}

	act, err := parseAction(msg.Signal)
	if err != nil {
		return nil, err
	}

	stg := parseStrategy(msg.Strategy)

	fmt.Printf("[%s] %s %s\n", s, stg, act)
	return &Signal{
		Strategy:  stg,
		Action:    act,
		Triggered: time.Now(),
	}, nil
}

func parseAction(msg string) (Action, error) {
	switch msg {
	case "空轉多訊號", "多方訊號":
		return LONG, nil
	case "多轉空訊號", "空方訊號":
		return SHORT, nil
	case "多方減倉訊號", "空方減倉訊號", "多方減倉50%", "多方減倉10%", "空方減倉10%", "空方減倉50%":
		return REDUCE, nil
	case "多方停損訊號", "空方停損訊號", "停損出場":
		return STOP_LOSS, nil
	case "多方平倉訊號", "多方平倉":
		return CLOSE_LONG, nil
	case "空方平倉訊號", "空方平倉":
		return CLOSE_SHORT, nil
	default:
		return UnknownSignal, fmt.Errorf("unknown alert:%v len:%d\n", msg, len(msg))
	}
}

func parseStrategy(msg string) Strategy {
	switch msg {
	case "左側拐點":
		return CounterTrading
	case "順勢指標", "順勢減倉":
		return TrendFollowing
	default:
		return UnknownStrategy
	}
}
