package structs

import (
	"time"
)

type HistoricalPrices struct {
	Success bool `json:"success"`
	Result  []struct {
		Close     float64   `json:"close"`
		High      float64   `json:"high"`
		Low       float64   `json:"low"`
		Open      float64   `json:"open"`
		StartTime time.Time `json:"startTime"`
		Volume    float64   `json:"volume"`
	} `json:"result"`
}

type Trades struct {
	Success bool `json:"success"`
	Result  []struct {
		ID          int64     `json:"id"`
		Liquidation bool      `json:"liquidation"`
		Price       float64   `json:"price"`
		Side        string    `json:"side"`
		Size        float64   `json:"size"`
		Time        time.Time `json:"time"`
	} `json:"result"`
}

type MarketResponse struct {
	Success bool   `json:"success"`
	Result  Market `json:"result"`
}

type Market struct {
	Name                  string  `json:"name"`
	Enabled               bool    `json:"enabled"`
	PostOnly              bool    `json:"postOnly"`
	PriceIncrement        float64 `json:"priceIncrement"`
	SizeIncrement         float64 `json:"sizeIncrement"`
	MinProvideSize        float64 `json:"minProvideSize"`
	Last                  float64 `json:"last"`
	Bid                   float64 `json:"bid"`
	Ask                   float64 `json:"ask"`
	Price                 float64 `json:"price"`
	Type                  string  `json:"type"`
	FutureType            *string `json:"futureType,omitempty"`
	BaseCurrency          *string `json:"baseCurrency,omitempty"`
	IsEtfMarket           bool    `json:"isEtfMarket"`
	QuoteCurrency         *string `json:"quoteCurrency,omitempty"`
	Underlying            *string `json:"underlying,omitempty"`
	Restricted            bool    `json:"restricted"`
	HighLeverageFeeExempt bool    `json:"highLeverageFeeExempt"`
	LargeOrderThreshold   float64 `json:"largeOrderThreshold"`
	Change1H              float64 `json:"change1h"`
	Change24H             float64 `json:"change24h"`
	ChangeBod             float64 `json:"changeBod"`
	QuoteVolume24H        float64 `json:"quoteVolume24h"`
	VolumeUsd24H          float64 `json:"volumeUsd24h"`
	PriceHigh24H          float64 `json:"priceHigh24h"`
	PriceLow24H           float64 `json:"priceLow24h"`
}

type OrderBookResponse struct {
	Success bool `json:"success"`
	Result  struct {
		Asks [][]float64 `json:"asks"`
		Bids [][]float64 `json:"bids"`
	} `json:"result"`
}
