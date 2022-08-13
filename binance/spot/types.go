package spot

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
	Symbols         []struct {
		Symbol                     string   `json:"symbol"`
		Status                     string   `json:"status"`
		BaseAsset                  string   `json:"baseAsset"`
		BaseAssetPrecision         int      `json:"baseAssetPrecision"`
		QuoteAsset                 string   `json:"quoteAsset"`
		QuotePrecision             int      `json:"quotePrecision"`
		QuoteAssetPrecision        int      `json:"quoteAssetPrecision"`
		BaseCommissionPrecision    int      `json:"baseCommissionPrecision"`
		QuoteCommissionPrecision   int      `json:"quoteCommissionPrecision"`
		OrderTypes                 []string `json:"orderTypes"`
		IcebergAllowed             bool     `json:"icebergAllowed"`
		OcoAllowed                 bool     `json:"ocoAllowed"`
		QuoteOrderQtyMarketAllowed bool     `json:"quoteOrderQtyMarketAllowed"`
		AllowTrailingStop          bool     `json:"allowTrailingStop"`
		CancelReplaceAllowed       bool     `json:"cancelReplaceAllowed"`
		IsSpotTradingAllowed       bool     `json:"isSpotTradingAllowed"`
		IsMarginTradingAllowed     bool     `json:"isMarginTradingAllowed"`
		Filters                    []struct {
			FilterType            string  `json:"filterType"`
			MinPrice              string  `json:"minPrice,omitempty"`
			MaxPrice              string  `json:"maxPrice,omitempty"`
			TickSize              float64 `json:"tickSize,string,omitempty"`
			MultiplierUp          string  `json:"multiplierUp,omitempty"`
			MultiplierDown        string  `json:"multiplierDown,omitempty"`
			AvgPriceMins          int     `json:"avgPriceMins,omitempty"`
			MinQty                string  `json:"minQty,omitempty"`
			MaxQty                string  `json:"maxQty,omitempty"`
			StepSize              string  `json:"stepSize,omitempty"`
			MinNotional           float64 `json:"minNotional,string,omitempty"`
			ApplyToMarket         bool    `json:"applyToMarket,omitempty"`
			Limit                 int     `json:"limit,omitempty"`
			MinTrailingAboveDelta int     `json:"minTrailingAboveDelta,omitempty"`
			MaxTrailingAboveDelta int     `json:"maxTrailingAboveDelta,omitempty"`
			MinTrailingBelowDelta int     `json:"minTrailingBelowDelta,omitempty"`
			MaxTrailingBelowDelta int     `json:"maxTrailingBelowDelta,omitempty"`
			MaxNumOrders          int     `json:"maxNumOrders,omitempty"`
			MaxNumAlgoOrders      int     `json:"maxNumAlgoOrders,omitempty"`
		} `json:"filters"`
		Permissions []string `json:"permissions"`
	} `json:"symbols"`
}

type AccountResp struct {
	MakerCommission  int64     `json:"makerCommission"`
	TakerCommission  int64     `json:"takerCommission"`
	BuyerCommission  int64     `json:"buyerCommission"`
	SellerCommission int64     `json:"sellerCommission"`
	CanTrade         bool      `json:"canTrade"`
	CanWithdraw      bool      `json:"canWithdraw"`
	CanDeposit       bool      `json:"canDeposit"`
	Balances         []Balance `json:"balances"`

	ErrorResp
}

type Balance struct {
	Asset  string  `json:"asset"`
	Free   float64 `json:"free,string"`
	Locked float64 `json:"locked,string"`
}

type MarginAccountResp struct {
	Assets              []MarginAccountAsset `json:"assets"`
	TotalAssetOfBtc     string               `json:"totalAssetOfBtc,omitempty"`
	TotalLiabilityOfBtc string               `json:"totalLiabilityOfBtc,omitempty"`
	TotalNetAssetOfBtc  string               `json:"totalNetAssetOfBtc,omitempty"`

	ErrorResp
}

type MarginAccountAsset struct {
	BaseAsset         Asset  `json:"baseAsset"`
	QuoteAsset        Asset  `json:"quoteAsset"`
	Symbol            string `json:"symbol"`
	IsolatedCreated   bool   `json:"isolatedCreated"`
	Enabled           bool   `json:"enabled"`
	MarginLevel       string `json:"marginLevel"`
	MarginLevelStatus string `json:"marginLevelStatus"`
	MarginRatio       string `json:"marginRatio"`
	IndexPrice        string `json:"indexPrice"`
	LiquidatePrice    string `json:"liquidatePrice"`
	LiquidateRate     string `json:"liquidateRate"`
	TradeEnabled      bool   `json:"tradeEnabled"`
}

type ErrorResp struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type Asset struct {
	Name          string `json:"asset"`
	BorrowEnabled bool   `json:"borrowEnabled"`
	Borrowed      string `json:"borrowed"`
	Free          string `json:"free"`
	Interest      string `json:"interest"`
	Locked        string `json:"locked"`
	Net           string `json:"netAsset"`
	NetOfBtc      string `json:"netAssetOfBtc"`
	RepayEnabled  bool   `json:"repayEnabled"`
	Total         string `json:"totalAsset"`
}

type MarginRepayResp struct {
	ErrorResp
	TxID      int    `json:"tranId"`
	ClientTag string `json:"clientTag"`
}

type OrderResp struct {
	ErrorResp
	Symbol                string `json:"symbol"`
	OrderId               int    `json:"orderId"`
	ClientOrderId         string `json:"clientOrderId"`
	TransactTime          int64  `json:"transactTime"`
	Price                 string `json:"price"`
	OrigQty               string `json:"origQty"`
	ExecutedQty           string `json:"executedQty"`
	CummulativeQuoteQty   string `json:"cummulativeQuoteQty"`
	Status                string `json:"status"`
	TimeInForce           string `json:"timeInForce"`
	Type                  string `json:"type"`
	Side                  string `json:"side"`
	MarginBuyBorrowAmount int    `json:"marginBuyBorrowAmount"`
	MarginBuyBorrowAsset  string `json:"marginBuyBorrowAsset"`
	IsIsolated            bool   `json:"isIsolated"`
	Fills                 []struct {
		Price           string `json:"price"`
		Qty             string `json:"qty"`
		Commission      string `json:"commission"`
		CommissionAsset string `json:"commissionAsset"`
	} `json:"fills"`
}

type MarginOrderResp struct {
	ClientOrderId       string `json:"clientOrderId"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	ExecutedQty         string `json:"executedQty"`
	IcebergQty          string `json:"icebergQty"`
	IsWorking           bool   `json:"isWorking"`
	OrderId             int    `json:"orderId"`
	OrigQty             string `json:"origQty"`
	Price               string `json:"price"`
	Side                string `json:"side"`
	Status              string `json:"status"`
	StopPrice           string `json:"stopPrice"`
	Symbol              string `json:"symbol"`
	IsIsolated          bool   `json:"isIsolated"`
	Time                int64  `json:"time"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	UpdateTime          int64  `json:"updateTime"`
}
