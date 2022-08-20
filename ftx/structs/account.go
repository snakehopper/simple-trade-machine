package structs

type AccountResponse struct {
	Success bool    `json:"success"`
	Result  Account `json:"result"`
}

type Account struct {
	BackstopProvider             bool       `json:"backstopProvider"`
	Collateral                   float64    `json:"collateral"`
	FreeCollateral               float64    `json:"freeCollateral"`
	InitialMarginRequirement     float64    `json:"initialMarginRequirement"`
	Leverage                     float64    `json:"leverage"`
	Liquidating                  bool       `json:"liquidating"`
	MaintenanceMarginRequirement float64    `json:"maintenanceMarginRequirement"`
	MakerFee                     float64    `json:"makerFee"`
	MarginFraction               float64    `json:"marginFraction"`
	OpenMarginFraction           float64    `json:"openMarginFraction"`
	TakerFee                     float64    `json:"takerFee"`
	TotalAccountValue            float64    `json:"totalAccountValue"`
	TotalPositionSize            float64    `json:"totalPositionSize"`
	Username                     string     `json:"username"`
	Positions                    []Position `json:"positions"`
}

type PositionResponse struct {
	Success bool       `json:"success"`
	Result  []Position `json:"result"`
}

type Position struct {
	Cost                         float64 `json:"cost"`
	EntryPrice                   float64 `json:"entryPrice"`
	EstimatedLiquidationPrice    float64 `json:"estimatedLiquidationPrice,omitempty"`
	Future                       string  `json:"future"`
	InitialMarginRequirement     float64 `json:"initialMarginRequirement"`
	LongOrderSize                float64 `json:"longOrderSize"`
	MaintenanceMarginRequirement float64 `json:"maintenanceMarginRequirement"`
	NetSize                      float64 `json:"netSize"`
	OpenSize                     float64 `json:"openSize"`
	RealizedPnl                  float64 `json:"realizedPnl"`
	ShortOrderSize               float64 `json:"shortOrderSize"`
	Side                         string  `json:"side"`
	Size                         float64 `json:"size"`
	UnrealizedPnl                float64 `json:"unrealizedPnl"`
}

type WalletBalancesResp struct {
	Success bool             `json:"success"`
	Result  []WalletBalances `json:"result"`
}

type WalletBalances struct {
	Coin                   string  `json:"coin"`
	Free                   float64 `json:"free"`
	SpotBorrow             float64 `json:"spotBorrow"`
	Total                  float64 `json:"total"`
	UsdValue               float64 `json:"usdValue"`
	AvailableWithoutBorrow float64 `json:"availableWithoutBorrow"`
}
