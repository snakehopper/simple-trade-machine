package ftx

import (
	"encoding/json"
	"fmt"
	"ghohoo.solutions/yt/ftx/structs"
	"strconv"
)

type NewOrder structs.NewOrder
type NewOrderResponse structs.NewOrderResponse
type OpenOrders structs.OpenOrders
type OrderHistory structs.OrderHistory
type NewTriggerOrder structs.NewTriggerOrder
type NewTriggerOrderResponse structs.NewTriggerOrderResponse
type OpenTriggerOrders structs.OpenTriggerOrders
type TriggerOrderHistory structs.TriggerOrderHistory
type Triggers structs.Triggers

func (a Api) GetOpenOrders(market string) (OpenOrders, error) {
	var openOrders OpenOrders
	resp, err := a._get("orders?market="+market, []byte(""))
	if err != nil {
		fmt.Println("Error GetOpenOrders", err)
		return openOrders, err
	}
	err = _processResponse(resp, &openOrders)
	return openOrders, err
}

func (a Api) GetOrderHistory(market string, startTime float64, endTime float64, limit int64) (OrderHistory, error) {
	var orderHistory OrderHistory
	requestBody, err := json.Marshal(map[string]interface{}{
		"market":     market,
		"start_time": startTime,
		"end_time":   endTime,
		"limit":      limit,
	})
	if err != nil {
		fmt.Println("Error GetOrderHistory", err)
		return orderHistory, err
	}
	resp, err := a._get("orders/history?market="+market, requestBody)
	if err != nil {
		fmt.Println("Error GetOrderHistory", err)
		return orderHistory, err
	}
	err = _processResponse(resp, &orderHistory)
	return orderHistory, err
}

func (a Api) GetOpenTriggerOrders(market string, _type string) (OpenTriggerOrders, error) {
	var openTriggerOrders OpenTriggerOrders
	requestBody, err := json.Marshal(map[string]string{"market": market, "type": _type})
	if err != nil {
		fmt.Println("Error GetOpenTriggerOrders", err)
		return openTriggerOrders, err
	}
	resp, err := a._get("conditional_orders?market="+market, requestBody)
	if err != nil {
		fmt.Println("Error GetOpenTriggerOrders", err)
		return openTriggerOrders, err
	}
	err = _processResponse(resp, &openTriggerOrders)
	return openTriggerOrders, err
}

func (a Api) GetTriggers(orderId string) (Triggers, error) {
	var trigger Triggers
	resp, err := a._get("conditional_orders/"+orderId+"/triggers", []byte(""))
	if err != nil {
		fmt.Println("Error GetTriggers", err)
		return trigger, err
	}
	err = _processResponse(resp, &trigger)
	return trigger, err
}

func (a Api) GetTriggerOrdersHistory(market string, startTime float64, endTime float64, limit int64) (TriggerOrderHistory, error) {
	var triggerOrderHistory TriggerOrderHistory
	requestBody, err := json.Marshal(map[string]interface{}{
		"market":     market,
		"start_time": startTime,
		"end_time":   endTime,
	})
	if err != nil {
		fmt.Println("Error GetTriggerOrdersHistory", err)
		return triggerOrderHistory, err
	}
	resp, err := a._get("conditional_orders/history?market="+market, requestBody)
	if err != nil {
		fmt.Println("Error GetTriggerOrdersHistory", err)
		return triggerOrderHistory, err
	}
	err = _processResponse(resp, &triggerOrderHistory)
	return triggerOrderHistory, err
}

func (a Api) PlaceOrder(market string, side string, price float64,
	_type string, size float64, reduceOnly bool, ioc bool, postOnly bool) (NewOrderResponse, error) {
	var newOrderResponse NewOrderResponse
	po := NewOrder{
		Market:     market,
		Side:       side,
		Price:      &price,
		Type:       _type,
		Size:       size,
		ReduceOnly: reduceOnly,
		Ioc:        ioc,
		PostOnly:   postOnly}
	if _type == "market" {
		po.Price = nil
	}
	requestBody, err := json.Marshal(po)
	if err != nil {
		fmt.Println("marshal post payload error:", err)
		return newOrderResponse, err
	}

	fmt.Println(">", string(requestBody))
	resp, err := a._post("orders", requestBody)
	if err != nil {
		fmt.Println("post PlaceOrder error:", err)
		return newOrderResponse, err
	}

	err = _processResponse(resp, &newOrderResponse)
	return newOrderResponse, err
}

func (a Api) PlaceTriggerOrder(market string, side string, size float64,
	_type string, reduceOnly bool, retryUntilFilled bool, triggerPrice float64,
	orderPrice float64, trailValue float64) (NewTriggerOrderResponse, error) {

	var newTriggerOrderResponse NewTriggerOrderResponse
	var newTriggerOrder NewTriggerOrder

	switch _type {
	case "stop":
		if orderPrice != 0 {
			newTriggerOrder = NewTriggerOrder{
				Market:       market,
				Side:         side,
				TriggerPrice: triggerPrice,
				Type:         _type,
				Size:         size,
				ReduceOnly:   reduceOnly,
				OrderPrice:   orderPrice,
			}
		} else {
			newTriggerOrder = NewTriggerOrder{
				Market:       market,
				Side:         side,
				TriggerPrice: triggerPrice,
				Type:         _type,
				Size:         size,
				ReduceOnly:   reduceOnly,
			}
		}
	case "trailingStop":
		newTriggerOrder = NewTriggerOrder{
			Market:     market,
			Side:       side,
			Type:       _type,
			Size:       size,
			ReduceOnly: reduceOnly,
			TrailValue: trailValue,
		}
	case "takeProfit":
		newTriggerOrder = NewTriggerOrder{
			Market:       market,
			Side:         side,
			TriggerPrice: triggerPrice,
			Type:         _type,
			Size:         size,
			ReduceOnly:   reduceOnly,
			OrderPrice:   orderPrice,
		}
	default:
		fmt.Println("Trigger type is not valid")
	}
	requestBody, err := json.Marshal(newTriggerOrder)
	if err != nil {
		fmt.Println("Error PlaceTriggerOrder", err)
		return newTriggerOrderResponse, err
	}
	resp, err := a._post("conditional_orders", requestBody)
	if err != nil {
		fmt.Println("Error PlaceTriggerOrder", err)
		return newTriggerOrderResponse, err
	}
	err = _processResponse(resp, &newTriggerOrderResponse)
	return newTriggerOrderResponse, err
}

func (a Api) CancelOrder(orderId int64) (Response, error) {
	var deleteResponse Response
	id := strconv.FormatInt(orderId, 10)
	resp, err := a._delete("orders/"+id, []byte(""))
	if err != nil {
		fmt.Println("Error CancelOrder", err)
		return deleteResponse, err
	}
	err = _processResponse(resp, &deleteResponse)
	return deleteResponse, err
}

func (a Api) CancelTriggerOrder(orderId int64) (Response, error) {
	var deleteResponse Response
	id := strconv.FormatInt(orderId, 10)
	resp, err := a._delete("conditional_orders/"+id, []byte(""))
	if err != nil {
		fmt.Println("Error CancelTriggerOrder", err)
		return deleteResponse, err
	}
	err = _processResponse(resp, &deleteResponse)
	return deleteResponse, err
}

func (a Api) CancelAllOrders() (Response, error) {
	var deleteResponse Response
	resp, err := a._delete("orders", []byte(""))
	if err != nil {
		fmt.Println("Error CancelAllOrders", err)
		return deleteResponse, err
	}
	err = _processResponse(resp, &deleteResponse)
	return deleteResponse, err
}
