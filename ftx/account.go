package ftx

import (
	"fmt"
	"ghohoo.solutions/yt/ftx/structs"
	"log"
)

type AccountInfo structs.AccountResponse
type Positions structs.PositionResponse

func (client *Client) GetAccount() (AccountInfo, error) {
	var acc AccountInfo
	resp, err := client._get("account", []byte(""))
	if err != nil {
		log.Printf("Error GetAccount: %v", err)
		return acc, err
	}
	err = _processResponse(resp, &acc)
	return acc, err
}

func (client *Client) GetPositions(showAvgPrice bool) (Positions, error) {
	link := "positions"
	if showAvgPrice {
		link = fmt.Sprintf("positions?showAvgPrice=true")
	}

	var positions Positions
	resp, err := client._get(link, []byte(""))
	if err != nil {
		log.Printf("Error GetPositions", err)
		return positions, err
	}
	err = _processResponse(resp, &positions)
	return positions, err
}
