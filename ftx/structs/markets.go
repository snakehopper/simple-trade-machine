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
	PostOnly              bool    `json:"post_only"`
	PriceIncrement        float64 `json:"price_increment"`
	SizeIncrement         float64 `json:"size_increment"`
	MinProvideSize        float64 `json:"min_provide_size"`
	Last                  float64 `json:"last"`
	Bid                   float64 `json:"bid"`
	Ask                   float64 `json:"ask"`
	Price                 float64 `json:"price"`
	Type                  string  `json:"type"`
	FutureType            string  `json:"future_type,omitempty"`
	BaseCurrency          string  `json:"base_currency,omitempty"`
	IsEtfMarket           bool    `json:"is_etf_market"`
	QuoteCurrency         string  `json:"quote_currency"`
	Underlying            string  `json:"underlying,omitempty"`
	Restricted            bool    `json:"restricted"`
	HighLeverageFeeExempt bool    `json:"high_leverage_fee_exempt"`
	LargeOrderThreshold   float64 `json:"large_order_threshold"`
	Change1h              float64 `json:"change_1_h"`
	Change24h             float64 `json:"change_24_h"`
	ChangeBod             float64 `json:"change_bod"`
	QuoteVolume24h        float64 `json:"quote_volume_24_h"`
	VolumeUsd24h          float64 `json:"volume_usd_24_h"`
	PriceHigh24h          float64 `json:"price_high_24_h"`
	PriceLow24h           float64 `json:"price_low_24_h"`
}
