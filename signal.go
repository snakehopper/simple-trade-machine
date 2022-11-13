//go:generate stringer -type Action -output signal_string.go
package function

import (
	"fmt"
	"regexp"
	"strconv"
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
	RAISE_LONG
	RAISE_SHORT

	CounterTrading  Strategy = "COUNTER"
	TrendFollowing  Strategy = "TREND"
	DoublePlay      Strategy = "DOUBLEPLAY"
	UnknownStrategy Strategy = "UnKNOwn"
)

type Action int
type Strategy string

type Signal struct {
	Strategy  Strategy
	Action    Action
	Triggered time.Time
	message   *AlertMessage
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
		message:   msg,
	}, nil
}

func (s Signal) ReducePct() (float64, error) {
	if s.Action != REDUCE {
		return 0, fmt.Errorf("not REDUCE action")
	}

	pct := regexp.MustCompilePOSIX(`([0-9]*[.])?[0-9]+`).FindString(s.message.Signal)
	return strconv.ParseFloat(pct, 64)
}

func (s Signal) RaisePct() (float64, error) {
	if s.Action != RAISE_SHORT && s.Action != RAISE_LONG {
		return 0, fmt.Errorf("not RAISE action")
	}

	pct := regexp.MustCompilePOSIX(`([0-9]*[.])?[0-9]+`).FindString(s.message.Signal)
	return strconv.ParseFloat(pct, 64)
}

func parseAction(msg string) (Action, error) {
	switch msg {
	case "空轉多訊號", "多方訊號":
		return LONG, nil
	case "多轉空訊號", "空方訊號":
		return SHORT, nil
	case "多方減倉訊號", "空方減倉訊號":
		return REDUCE, nil
	case "多方停損訊號", "空方停損訊號", "停損出場":
		return STOP_LOSS, nil
	case "多方平倉訊號", "多方平倉":
		return CLOSE_LONG, nil
	case "空方平倉訊號", "空方平倉":
		return CLOSE_SHORT, nil
	default:
		//retry message with variable
		switch {
		case strings.HasPrefix(msg, "多方減倉"), strings.HasPrefix(msg, "空方減倉"):
			return REDUCE, nil
		case strings.HasPrefix(msg, "多方加倉"):
			return RAISE_LONG, nil
		case strings.HasPrefix(msg, "空方加倉"):
			return RAISE_SHORT, nil
		}
		return UnknownSignal, fmt.Errorf("unknown alert:%v len:%d\n", msg, len(msg))
	}
}

func parseStrategy(msg string) Strategy {
	switch msg {
	case "左側拐點":
		return CounterTrading
	case "順勢指標", "順勢減倉":
		return TrendFollowing
	case "兩手策略":
		return DoublePlay
	default:
		return UnknownStrategy
	}
}
