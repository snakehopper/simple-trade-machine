package data

type Exchange interface {
	//MaxQuoteValue return order available to quote, e.g. free collateral * leverage
	//Spot instrument might return without leverage
	//sym provided for isolated margin account
	MaxQuoteValue(sym string) (total, free float64, err error)

	//GetPair return trading pair static info, which always can be cached
	GetPair(sym string) (*Pair, error)

	GetMarket(sym string) (*Market, error)

	GetPosition(sym string) (float64, error)
	LimitOrder(sym string, side Side, px float64, qty float64, ioc bool, postOnly bool) error
	MarketOrder(sym string, side Side, px *float64, qty *float64) error
}

type Side string

const (
	Buy  Side = "buy"
	Sell Side = "sell"
)

type Market struct {
	Type        MarketType
	Bid         float64
	Ask         float64
	Last        float64
	TickSize    float64
	MinNotional float64
}

type MarketType string

var (
	Spot   MarketType = "spot"
	Future MarketType = "future"
)

type Pair struct {
	Name  string
	Type  MarketType
	Base  string
	Quote string
}

func (p Pair) IsFuture() bool {
	return p.Type == Future
}

func (p Pair) IsSpot() bool {
	return p.Type == Spot
}
