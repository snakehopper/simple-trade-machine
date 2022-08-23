package function

import (
	"fmt"
	"strings"
)

const (
	UnknownSignal Signal = iota
	LONG
	SHORT
	REDUCE
	CLOSE
	STOP_LOSS
)

type Signal int

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

func NewSignal(s string) Signal {
	if len(s) == 0 {
		return UnknownSignal
	}

	msg, err := NewAlertMessage(s)
	if err != nil {
		fmt.Println(msg)
		return UnknownSignal
	}

	switch msg.Signal {
	case "空轉多訊號", "多方訊號":
		return LONG
	case "多轉空訊號", "空方訊號":
		return SHORT
	case "多方減倉訊號", "空方減倉訊號":
		return REDUCE
	case "多方停損訊號", "空方停損訊號":
		return STOP_LOSS
	case "多方平倉訊號", "空方平倉訊號":
		return CLOSE
	default:
		fmt.Printf("unknown alert:%v len:%d\n", s, len(s))
		return UnknownSignal
	}
}
