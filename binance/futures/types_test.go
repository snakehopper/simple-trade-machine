package futures

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApi_OrderBook(t *testing.T) {
	data := `{
  "lastUpdateId": 1027024,
  "bids": [
    [
      "4.00000000",
      "431.00000000"
    ]
  ],
  "asks": [
    [
      "4.00000200",
      "12.00000000"
    ]
  ]
}`
	var res OrderBookResp
	err := json.Unmarshal([]byte(data), &res)
	assert.Nil(t, err)
	fmt.Println(err)
	for _, bid := range res.Bids {
		px, err := bid[0].Float64()
		assert.Nil(t, err)
		sz, err := bid[1].Float64()
		assert.Nil(t, err)
		fmt.Println(px, sz)
	}
}
