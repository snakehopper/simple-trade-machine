package function

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAlert(t *testing.T) {
	dat := []struct {
		input  string
		output Action
	}{
		{"左側拐點｜多方訊號｜1小時｜$0.00001188", LONG},
		{"順勢指標｜多轉空訊號｜30分鐘｜$23322.13", SHORT},
		{"順勢指標｜空轉多訊號｜30分鐘｜$23322.13", LONG},
		{"順勢指標｜多方平倉訊號｜30分鐘｜$23322.13", CLOSE},
		{"順勢指標｜空方平倉訊號｜30分鐘｜$23322.13", CLOSE},
		{"左側拐點｜多方平倉訊號｜1小時｜$0.000012", CLOSE},
		{"順勢指標｜多方減倉訊號｜30分鐘｜$23371.23", REDUCE},
		{"順勢指標｜空方減倉訊號｜30分鐘｜$23371.23", REDUCE},
		{"左側拐點｜多方停損訊號｜1小時｜$0.0000117802", STOP_LOSS},
		{"左側拐點｜空方停損訊號｜1小時｜$0.0000117802", STOP_LOSS},
		{"左側拐點｜停損出場｜｜$19589.7279166667", STOP_LOSS},
		{"順勢減倉｜空方平倉｜｜$1547.3895294034", CLOSE_SHORT},
		{"順勢減倉｜多方平倉｜｜$1547.3895294034", CLOSE_LONG},
	}

	for _, d := range dat {
		sig, err := NewSignal(d.input)
		assert.Nil(t, err)

		assert.Equal(t, d.output, sig.Action)
	}
}
