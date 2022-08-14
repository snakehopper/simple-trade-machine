package futures

import "fmt"

type ExchangeInfoResp struct {
	Timezone   string `json:"timezone"`
	ServerTime int64  `json:"serverTime"`
	RateLimits []struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		IntervalNum   int    `json:"intervalNum"`
		Limit         int    `json:"limit"`
	} `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Symbols         []SymbolResp  `json:"symbols"`
}

type SymbolResp struct {
	Symbol                string   `json:"symbol"`
	Pair                  string   `json:"pair"`
	ContractType          string   `json:"contractType"`
	DeliveryDate          int64    `json:"deliveryDate"`
	OnboardDate           int64    `json:"onboardDate"`
	Status                string   `json:"status"`
	MaintMarginPercent    string   `json:"maintMarginPercent"`
	RequiredMarginPercent string   `json:"requiredMarginPercent"`
	BaseAsset             string   `json:"baseAsset"`
	QuoteAsset            string   `json:"quoteAsset"`
	MarginAsset           string   `json:"marginAsset"`
	PricePrecision        int      `json:"pricePrecision"`
	QuantityPrecision     int      `json:"quantityPrecision"`
	BaseAssetPrecision    int      `json:"baseAssetPrecision"`
	QuotePrecision        int      `json:"quotePrecision"`
	UnderlyingType        string   `json:"underlyingType"`
	UnderlyingSubType     []string `json:"underlyingSubType"`
	SettlePlan            int      `json:"settlePlan"`
	TriggerProtect        string   `json:"triggerProtect"`
	Filters               []struct {
		FilterType        string  `json:"filterType"`
		MaxPrice          string  `json:"maxPrice,omitempty"`
		MinPrice          string  `json:"minPrice,omitempty"`
		TickSize          float64 `json:"tickSize,string,omitempty"`
		MaxQty            string  `json:"maxQty,omitempty"`
		MinQty            string  `json:"minQty,omitempty"`
		StepSize          string  `json:"stepSize,omitempty"`
		Limit             int     `json:"limit,omitempty"`
		Notional          float64 `json:"notional,string,omitempty"`
		MultiplierUp      string  `json:"multiplierUp,omitempty"`
		MultiplierDown    string  `json:"multiplierDown,omitempty"`
		MultiplierDecimal int     `json:"multiplierDecimal,string,omitempty"`
	} `json:"filters"`
	OrderType       []string `json:"OrderType"`
	TimeInForce     []string `json:"timeInForce"`
	LiquidationFee  string   `json:"liquidationFee"`
	MarketTakeBound string   `json:"marketTakeBound"`
}

type AccountResp struct {
	FeeTier                     int        `json:"feeTier"`
	CanTrade                    bool       `json:"canTrade"`
	CanDeposit                  bool       `json:"canDeposit"`
	CanWithdraw                 bool       `json:"canWithdraw"`
	UpdateTime                  int        `json:"updateTime"`
	TotalInitialMargin          string     `json:"totalInitialMargin"`
	TotalMaintMargin            string     `json:"totalMaintMargin"`
	TotalWalletBalance          string     `json:"totalWalletBalance"`
	TotalUnrealizedProfit       string     `json:"totalUnrealizedProfit"`
	TotalMarginBalance          float64    `json:"totalMarginBalance,string"`
	TotalPositionInitialMargin  string     `json:"totalPositionInitialMargin"`
	TotalOpenOrderInitialMargin string     `json:"totalOpenOrderInitialMargin"`
	TotalCrossWalletBalance     string     `json:"totalCrossWalletBalance"`
	TotalCrossUnPnl             string     `json:"totalCrossUnPnl"`
	AvailableBalance            float64    `json:"availableBalance,string"`
	MaxWithdrawAmount           float64    `json:"maxWithdrawAmount,string"`
	Assets                      []Asset    `json:"assets"`
	Positions                   []Position `json:"positions"`

	ErrorResp
}

type Asset struct {
	Asset                  string `json:"asset"`
	WalletBalance          string `json:"walletBalance"`
	UnrealizedProfit       string `json:"unrealizedProfit"`
	MarginBalance          string `json:"marginBalance"`
	MaintMargin            string `json:"maintMargin"`
	InitialMargin          string `json:"initialMargin"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	CrossWalletBalance     string `json:"crossWalletBalance"`
	CrossUnPnl             string `json:"crossUnPnl"`
	AvailableBalance       string `json:"availableBalance"`
	MaxWithdrawAmount      string `json:"maxWithdrawAmount"`
	MarginAvailable        bool   `json:"marginAvailable"`
	UpdateTime             int64  `json:"updateTime"`
}

type Position struct {
	Symbol                 string  `json:"symbol"`
	InitialMargin          string  `json:"initialMargin"`
	MaintMargin            string  `json:"maintMargin"`
	UnrealizedProfit       string  `json:"unrealizedProfit"`
	PositionInitialMargin  string  `json:"positionInitialMargin"`
	OpenOrderInitialMargin string  `json:"openOrderInitialMargin"`
	Leverage               float64 `json:"leverage,string"`
	Isolated               bool    `json:"isolated"`
	EntryPrice             string  `json:"entryPrice"`
	MaxNotional            string  `json:"maxNotional"`
	BidNotional            string  `json:"bidNotional"`
	AskNotional            string  `json:"askNotional"`
	PositionSide           string  `json:"positionSide"`
	PositionAmt            float64 `json:"positionAmt,string"`
	UpdateTime             int     `json:"updateTime"`
}

func (a AccountResp) GetPosition(sym string) (Position, error) {
	for _, pos := range a.Positions {
		if pos.Symbol == sym {
			return pos, nil
		}
	}

	return Position{}, fmt.Errorf("%s Position not found", sym)
}

type ErrorResp struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type OrderResp struct {
	ClientOrderId string `json:"clientOrderId"`
	CumQty        string `json:"cumQty"`
	CumQuote      string `json:"cumQuote"`
	ExecutedQty   string `json:"executedQty"`
	OrderId       int    `json:"orderId"`
	AvgPrice      string `json:"avgPrice"`
	OrigQty       string `json:"origQty"`
	Price         string `json:"price"`
	ReduceOnly    bool   `json:"reduceOnly"`
	Side          string `json:"side"`
	PositionSide  string `json:"positionSide"`
	Status        string `json:"status"`
	StopPrice     string `json:"stopPrice"`
	ClosePosition bool   `json:"closePosition"`
	Symbol        string `json:"symbol"`
	TimeInForce   string `json:"timeInForce"`
	Type          string `json:"type"`
	OrigType      string `json:"origType"`
	ActivatePrice string `json:"activatePrice"`
	PriceRate     string `json:"priceRate"`
	UpdateTime    int64  `json:"updateTime"`
	WorkingType   string `json:"workingType"`
	PriceProtect  bool   `json:"priceProtect"`

	ErrorResp
}
