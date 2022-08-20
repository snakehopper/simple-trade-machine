package ftx

import (
	"fmt"
	"ghohoo.solutions/yt/ftx/structs"
)

type AccountInfo structs.AccountResponse
type Positions structs.PositionResponse

func (a Api) GetAccount() (AccountInfo, error) {
	var acc AccountInfo
	resp, err := a._get("account", []byte(""))
	if err != nil {
		fmt.Printf("Error GetAccount: %v\n", err)
		return acc, err
	}
	err = _processResponse(resp, &acc)
	return acc, err
}

func (a Api) GetPositions(showAvgPrice bool) (Positions, error) {
	link := "positions"
	if showAvgPrice {
		link = fmt.Sprintf("positions?showAvgPrice=true")
	}

	var positions Positions
	resp, err := a._get(link, []byte(""))
	if err != nil {
		fmt.Printf("Error GetPositions: %v\n", err)
		return positions, err
	}
	err = _processResponse(resp, &positions)
	return positions, err
}

func (a Api) GetBalance(sym string) (*structs.WalletBalances, error) {
	var balances structs.WalletBalancesResp
	resp, err := a._get("wallet/balances", []byte(""))
	if err != nil {
		fmt.Printf("Error GetAccount: %v\n", err)
		return nil, err
	}
	if err = _processResponse(resp, &balances); err != nil {
		return nil, err
	}

	pair := a.GetTradingPair(sym)
	for _, res := range balances.Result {
		if res.Coin == pair.Base {
			return &res, nil
		}
	}
	return nil, fmt.Errorf("invalid symbol: %v", sym)
}
