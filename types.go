package function

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
	LongPosition(sym string) error
	ReducePosition(sym string) error
	ClosePosition(sym string) error
	StopLossPosition(sym string) error
}
