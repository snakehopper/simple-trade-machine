package ftx

import (
	"encoding/json"
	"fmt"
	"ghohoo.solutions/yt/ftx/structs"
)

type SubaccountsList structs.SubaccountsList
type Subaccount structs.Subaccount
type Response structs.Response
type SubaccountBalances structs.SubaccountBalances
type TransferSubaccounts structs.TransferSubaccounts

func (a Api) GetSubaccounts() (SubaccountsList, error) {
	var subaccounts SubaccountsList
	resp, err := a._get("subaccounts", []byte(""))
	if err != nil {
		fmt.Printf("Error GetSubaccounts: %v\n", err)
		return subaccounts, err
	}
	err = _processResponse(resp, &subaccounts)
	return subaccounts, err
}

func (a Api) CreateSubaccount(nickname string) (Subaccount, error) {
	var subaccount Subaccount
	requestBody, err := json.Marshal(map[string]string{"nickname": nickname})
	if err != nil {
		fmt.Printf("Error CreateSubaccount: %v\n", err)
		return subaccount, err
	}
	resp, err := a._post("subaccounts", requestBody)
	if err != nil {
		fmt.Printf("Error CreateSubaccount: %v\n", err)
		return subaccount, err
	}
	err = _processResponse(resp, &subaccount)
	return subaccount, err
}

func (a Api) ChangeSubaccountName(nickname string, newNickname string) (Response, error) {
	var changeSubaccount Response
	requestBody, err := json.Marshal(map[string]string{"nickname": nickname, "newNickname": newNickname})
	if err != nil {
		fmt.Printf("Error ChangeSubaccountName: %v\n", err)
		return changeSubaccount, err
	}
	resp, err := a._post("subaccounts/update_name", requestBody)
	if err != nil {
		fmt.Printf("Error ChangeSubaccountName: %v\n", err)
		return changeSubaccount, err
	}
	err = _processResponse(resp, &changeSubaccount)
	return changeSubaccount, err
}

func (a Api) DeleteSubaccount(nickname string) (Response, error) {
	var deleteSubaccount Response
	requestBody, err := json.Marshal(map[string]string{"nickname": nickname})
	if err != nil {
		fmt.Printf("Error DeleteSubaccount: %v\n", err)
		return deleteSubaccount, err
	}
	resp, err := a._delete("subaccounts", requestBody)
	if err != nil {
		fmt.Printf("Error DeleteSubaccount: %v\n", err)
		return deleteSubaccount, err
	}
	err = _processResponse(resp, &deleteSubaccount)
	return deleteSubaccount, err
}

func (a Api) GetSubaccountBalances(nickname string) (SubaccountBalances, error) {
	var subaccountBalances SubaccountBalances
	resp, err := a._get("subaccounts/"+nickname+"/balances", []byte(""))
	if err != nil {
		fmt.Printf("Error SubaccountBalances: %v\n", err)
		return subaccountBalances, err
	}
	err = _processResponse(resp, &subaccountBalances)
	return subaccountBalances, err
}

func (a Api) TransferSubaccounts(coin string, size float64, source string, destination string) (TransferSubaccounts, error) {
	var transferSubaccounts TransferSubaccounts
	requestBody, err := json.Marshal(map[string]interface{}{
		"coin":        coin,
		"size":        size,
		"source":      source,
		"destination": destination,
	})
	if err != nil {
		fmt.Printf("Error TransferSubaccounts: %v\n", err)
		return transferSubaccounts, err
	}
	resp, err := a._post("subaccounts/transfer", requestBody)
	if err != nil {
		fmt.Printf("Error TransferSubaccounts: %v\n", err)
		return transferSubaccounts, err
	}
	err = _processResponse(resp, &transferSubaccounts)
	return transferSubaccounts, err
}
