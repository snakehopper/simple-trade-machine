package data

type Exchange interface {
	//MaxQuoteValue return order available to quote, e.g. free collateral * leverage
	//Spot instrument might return without leverage
	//sym provided for isolated margin account
	MaxQuoteValue(sym string) (total, free float64, err error)

	//GetPair return trading pair static info, which always can be cached
	GetPair(sym string) (*Pair, error)

	GetMarket(sym string) (*Market, error)
	GetOrderBook(sym string) (*OrderBook, error)

	//GetPosition return signed position
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
	TickSize    float64
	MinNotional float64
}

type MarketType string

var (
	Spot   MarketType = "spot"
	Future MarketType = "future"
)

type OrderType string

var (
	LimitOrder  OrderType = "limit"
	MarketOrder OrderType = "market"
)

type OrderBook struct {
	Bid []OrderBookLevel
	Ask []OrderBookLevel
}

func (b OrderBook) MidPx() float64 {
	return (b.Ask[0].Px + b.Bid[0].Px) / 2
}

type OrderBookLevel struct {
	Px   float64
	Size float64
}

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
