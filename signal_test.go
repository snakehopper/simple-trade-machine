package function

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSignal(t *testing.T) {
	dat := []struct {
		input  string
		output Signal
	}{
		{"左側拐點｜多方訊號｜1小時｜$0.00001188", Signal{Strategy: CounterTrading, Action: LONG}},
		{"順勢指標｜多轉空訊號｜30分鐘｜$23322.13", Signal{Strategy: TrendFollowing, Action: SHORT}},
		{"順勢指標｜空轉多訊號｜30分鐘｜$23322.13", Signal{Strategy: TrendFollowing, Action: LONG}},
		{"順勢指標｜多方平倉訊號｜30分鐘｜$23322.13", Signal{Strategy: TrendFollowing, Action: CLOSE_LONG}},
		{"順勢指標｜空方平倉訊號｜30分鐘｜$23322.13", Signal{Strategy: TrendFollowing, Action: CLOSE_SHORT}},
		{"左側拐點｜多方平倉訊號｜1小時｜$0.000012", Signal{Strategy: CounterTrading, Action: CLOSE_LONG}},
		{"順勢指標｜多方減倉訊號｜30分鐘｜$23371.23", Signal{Strategy: TrendFollowing, Action: REDUCE}},
		{"順勢指標｜空方減倉訊號｜30分鐘｜$23371.23", Signal{Strategy: TrendFollowing, Action: REDUCE}},
		{"左側拐點｜多方停損訊號｜1小時｜$0.0000117802", Signal{Strategy: CounterTrading, Action: STOP_LOSS}},
		{"左側拐點｜停損出場｜｜$19589.7279166667", Signal{Strategy: CounterTrading, Action: STOP_LOSS}},
		{"順勢減倉｜空方平倉｜｜$1547.3895294034", Signal{Strategy: TrendFollowing, Action: CLOSE_SHORT}},
		{"順勢減倉｜多方平倉｜｜$1547.3895294034", Signal{Strategy: TrendFollowing, Action: CLOSE_LONG}},
		{input: "左側拐點｜多方訊號｜｜$4.747", output: Signal{Strategy: CounterTrading, Action: LONG}},
		{input: "左側拐點｜多方減倉50%｜2小時｜$0.00436", output: Signal{Strategy: CounterTrading, Action: REDUCE}},
		{input: "左側拐點｜停損出場｜｜$19680.7266527778", output: Signal{Strategy: CounterTrading, Action: STOP_LOSS}},
		{input: "左側拐點｜多方平倉｜｜$19774.8945274472", output: Signal{Strategy: CounterTrading, Action: CLOSE_LONG}},
		{input: "順勢減倉｜多方訊號｜3小時｜$0.69015", output: Signal{Strategy: TrendFollowing, Action: LONG}},
		{input: "順勢減倉｜空方訊號｜｜$1550.9487759082", output: Signal{Strategy: TrendFollowing, Action: SHORT}},
		{input: "順勢減倉｜空方減倉10%｜1天｜$14.035", output: Signal{Strategy: TrendFollowing, Action: REDUCE}},
		{input: "順勢減倉｜多方減倉10%｜1天｜$14.035", output: Signal{Strategy: TrendFollowing, Action: REDUCE}},
		{input: "順勢減倉｜空方平倉｜3小時｜$0.69015", output: Signal{Strategy: TrendFollowing, Action: CLOSE_SHORT}},
		{input: "順勢減倉｜多方平倉｜3小時｜$0.69015", output: Signal{Strategy: TrendFollowing, Action: CLOSE_LONG}},

		{input: "順勢減倉｜空方減倉10%｜1分鐘｜$2.0295", output: Signal{Strategy: TrendFollowing, Action: REDUCE}},
		{input: "順勢減倉｜多方訊號｜1分鐘｜$2.0714", output: Signal{Strategy: TrendFollowing, Action: LONG}},
		{input: "左側拐點｜多方減倉50%｜1分鐘｜$283.4", output: Signal{Strategy: CounterTrading, Action: REDUCE}},
		{input: "左側拐點｜停損出場｜1分鐘｜$285.1960388889", output: Signal{Strategy: CounterTrading, Action: STOP_LOSS}},
		{input: "左側拐點｜多方訊號｜1分鐘｜$286.3", output: Signal{Strategy: CounterTrading, Action: LONG}},
		{input: "兩手策略｜多方平倉｜1分鐘｜$16838.5203060912", output: Signal{Strategy: DoublePlay, Action: CLOSE_LONG}},
		{input: "兩手策略｜多方減倉10%｜1分鐘｜$16894.9124278729", output: Signal{Strategy: DoublePlay, Action: REDUCE}},
		{input: "兩手策略｜多方加倉10%｜1分鐘｜$16909.79", output: Signal{Strategy: DoublePlay, Action: RAISE_LONG}},
		{input: "兩手策略｜多方減倉20%｜1分鐘｜$16890.5602685398", output: Signal{Strategy: DoublePlay, Action: REDUCE}},
		{input: "兩手策略｜多方訊號｜1分鐘｜$16875.39", output: Signal{Strategy: DoublePlay, Action: LONG}},
	}

	for _, d := range dat {
		sig, err := NewSignal(d.input)
		assert.Nil(t, err)

		assert.Equal(t, d.output.Strategy, sig.Strategy, d.input)
		assert.Equal(t, d.output.Action, sig.Action, d.input)
	}
}

func TestNewSignal_reduce_percent(t *testing.T) {
	dat := []struct {
		input       string
		output      float64
		expectedErr bool
	}{
		{input: "左側拐點｜多方減倉50%｜2小時｜$0.00436", output: 50},
		{input: "順勢減倉｜空方減倉10%｜1天｜$14.035", output: 10},
		{input: "順勢減倉｜多方減倉12.3%｜1天｜$14.035", output: 12.3},
		{input: "順勢減倉｜空方減倉10%｜1分鐘｜$2.0295", output: 10},
		{input: "順勢減倉｜多方訊號｜1分鐘｜$2.0714", output: 0, expectedErr: true},
		{input: "左側拐點｜多方減倉50%｜1分鐘｜$283.4", output: 50},
		{input: "左側拐點｜停損出場｜1分鐘｜$285.1960388889", output: 0, expectedErr: true},
		{input: "左側拐點｜多方訊號｜1分鐘｜$286.3", output: 0, expectedErr: true},
		{input: "兩手策略｜多方平倉｜1分鐘｜$16838.5203060912", output: 0, expectedErr: true},
		{input: "兩手策略｜多方減倉10%｜1分鐘｜$16894.9124278729", output: 10},
		{input: "兩手策略｜多方加倉10%｜1分鐘｜$16909.79", output: 0, expectedErr: true},
		{input: "兩手策略｜多方減倉20%｜1分鐘｜$16890.5602685398", output: 20},
		{input: "兩手策略｜多方訊號｜1分鐘｜$16875.39", output: 0, expectedErr: true},
	}

	for _, d := range dat {
		sig, err := NewSignal(d.input)
		assert.Nil(t, err)

		pct, err := sig.ReducePct()
		if d.expectedErr {
			assert.NotNil(t, d.expectedErr)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, pct, d.output)
		}
	}

	t.Run("not reduce action", func(t *testing.T) {
		sig, err := NewSignal("順勢減倉｜多方訊號｜3小時｜$0.69015")
		assert.Nil(t, err)
		_, err = sig.ReducePct()
		assert.NotNil(t, err)
	})

}

func TestNewSignal_raise_percent(t *testing.T) {
	dat := []struct {
		input       string
		output      float64
		expectedErr bool
	}{
		{input: "左側拐點｜多方減倉50%｜2小時｜$0.00436", output: 0, expectedErr: true},
		{input: "順勢減倉｜空方減倉10%｜1天｜$14.035", output: 0, expectedErr: true},
		{input: "順勢減倉｜多方減倉12.3%｜1天｜$14.035", output: 0, expectedErr: true},
		{input: "順勢減倉｜多方訊號｜1分鐘｜$2.0714", output: 0, expectedErr: true},
		{input: "左側拐點｜多方減倉50%｜1分鐘｜$283.4", output: 0, expectedErr: true},
		{input: "左側拐點｜停損出場｜1分鐘｜$285.1960388889", output: 0, expectedErr: true},
		{input: "左側拐點｜多方訊號｜1分鐘｜$286.3", output: 0, expectedErr: true},
		{input: "兩手策略｜多方平倉｜1分鐘｜$16838.5203060912", output: 0, expectedErr: true},
		{input: "兩手策略｜多方減倉10%｜1分鐘｜$16894.9124278729", output: 0, expectedErr: true},
		{input: "兩手策略｜多方加倉10%｜1分鐘｜$16909.79", output: 10},
		{input: "兩手策略｜空方加倉12.3%｜1分鐘｜$16909.79", output: 12.3},
	}

	for _, d := range dat {
		sig, err := NewSignal(d.input)
		assert.Nil(t, err)

		pct, err := sig.RaisePct()
		if d.expectedErr {
			assert.NotNil(t, d.expectedErr)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, pct, d.output)
		}
	}

	t.Run("not raise action", func(t *testing.T) {
		sig, err := NewSignal("順勢減倉｜多方減倉10%｜3小時｜$0.69015")
		assert.Nil(t, err)
		_, err = sig.RaisePct()
		assert.NotNil(t, err)
	})

}
