package function

import (
	"ghohoo.solutions/yt/internal/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAlert(t *testing.T) {
	dat := []struct {
		input  string
		output data.Alert
	}{
		{"左側拐點｜多方訊號｜1小時｜$0.00001188", data.LONG},
		{"順勢指標｜多轉空訊號｜30分鐘｜$23322.13", data.SHORT},
		{"順勢指標｜空轉多訊號｜30分鐘｜$23322.13", data.LONG},
		{"順勢指標｜多方平倉訊號｜30分鐘｜$23322.13", data.CLOSE},
		{"順勢指標｜空方平倉訊號｜30分鐘｜$23322.13", data.CLOSE},
		{"左側拐點｜多方平倉訊號｜1小時｜$0.000012", data.CLOSE},
		{"順勢指標｜多方減倉訊號｜30分鐘｜$23371.23", data.REDUCE},
		{"順勢指標｜空方減倉訊號｜30分鐘｜$23371.23", data.REDUCE},
		{"左側拐點｜多方停損訊號｜1小時｜$0.0000117802", data.STOP_LOSS},
		{"左側拐點｜空方停損訊號｜1小時｜$0.0000117802", data.STOP_LOSS},
	}

	for _, d := range dat {
		assert.Equal(t, d.output, data.NewAlert(d.input))
	}
}
