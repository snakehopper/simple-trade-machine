package spot

import (
	"encoding/json"
	"fmt"
	"ghohoo.solutions/yt/binance/com"
	"ghohoo.solutions/yt/internal/data"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Api struct {
	*com.Client
	log *zap.SugaredLogger
}

func NewApi(logger *zap.SugaredLogger, apiKey, secret string) *Api {
	return &Api{
		Client: com.NewClient(logger, "https://api.binance.com", apiKey, secret),
		log:    logger,
	}
}

func (a Api) MaxQuoteValue(sym string) (total, free float64, err error) {
	acc, err := a.AccountInfo()
	if err != nil {
		return
	} else if acc.Code != 0 {
		err = fmt.Errorf("fetch account error code:%v msg:%v", acc.Code, acc.Msg)
		return
	}

	mx := viper.GetFloat64("SPOT_OPEN_X")
	quote := a.GetTradingPair(sym).Quote
	for _, bal := range acc.Balances {
		if bal.Asset == quote {
			free = bal.Free * mx
			total = (bal.Free + bal.Locked) * mx
			return
		}
	}
	return
}

func (a Api) GetMarket(sym string) (*data.Market, error) {
	res, err := a.ExchangeInfo(sym)
	if err != nil {
		return nil, err
	}

	for _, rs := range res.Symbols {
		if rs.Symbol == sym {
			var tickSize, minNotional float64
			for _, fp := range rs.Filters {
				if fp.FilterType == "MIN_NOTIONAL" {
					minNotional = fp.MinNotional
				} else if fp.FilterType == "PRICE_FILTER" {
					tickSize = fp.TickSize
				}
			}
			return &data.Market{
				Type:        data.Spot,
				TickSize:    tickSize,
				MinNotional: minNotional,
			}, nil
		}
	}
	return nil, fmt.Errorf("symbol %v not found", sym)
}

func (a Api) GetOrderBook(sym string) (*data.OrderBook, error) {
	res, err := a.OrderBook(sym)
	if err != nil {
		return nil, err
	}

	var out = &data.OrderBook{
		Bid: make([]data.OrderBookLevel, 0),
		Ask: make([]data.OrderBookLevel, 0),
	}
	for _, ask := range res.Asks {
		px, _ := ask[0].Float64()
		sz, _ := ask[1].Float64()
		out.Ask = append(out.Ask, data.OrderBookLevel{Px: px, Size: sz})
	}
	for _, bid := range res.Bids {
		px, _ := bid[0].Float64()
		sz, _ := bid[1].Float64()
		out.Bid = append(out.Ask, data.OrderBookLevel{Px: px, Size: sz})
	}

	return out, nil
}

func (a Api) GetPair(sym string) (*data.Pair, error) {
	p := a.GetTradingPair(sym)
	if p.Name == "" {
		return nil, fmt.Errorf("invalid symbol")
	}

	return &p, nil
}

func (a Api) GetPosition(sym string) (float64, error) {
	acc, err := a.AccountInfo()
	if err != nil {
		return 0, err
	} else if acc.Code != 0 {
		err = fmt.Errorf("fetch account error code:%v msg:%v", acc.Code, acc.Msg)
		return 0, err
	}

	exch, err := a.ExchangeInfo(sym)
	if err != nil {
		return 0, err
	}
	if len(exch.Symbols) == 0 {
		return 0, fmt.Errorf("unkonwn symbol:%s", sym)
	}

	var base string
	if s := exch.Symbols[0]; s.Symbol != sym {
		err = fmt.Errorf("exchange info symbol not match: %+v", s)
		return 0, err
	} else {
		base = s.BaseAsset
	}
	for _, bal := range acc.Balances {
		if bal.Asset == base {
			return bal.Free, nil
		}
	}

	return 0, fmt.Errorf("unknown symbol %v", sym)
}

func (a Api) LimitOrder(sym string, side data.Side, px float64, qty float64, ioc, _postOnly, _reduceOnly bool) (string, error) {
	exch, err := a.ExchangeInfo(sym)
	if err != nil {
		return "", err
	}

	var v = url.Values{}
	v.Set("side", strings.ToUpper(string(side)))
	v.Set("symbol", sym)
	v.Set("type", "LIMIT")
	if ioc {
		v.Set("timeInForce", "IOC")
	} else {
		v.Set("timeInForce", "GTC")
	}
	rounded := exch.RoundLotSize(sym, math.Abs(qty))
	if rounded == 0 {
		a.log.Infof("%f rounded quantity is 0, skip place order", qty)
		return "", nil
	} else if rounded < exch.MinNotional(sym) {
		a.log.Info("rounded quoteOrderQty smaller than MIN_NOTIONAL, skip place order")
		return "", nil
	}

	v.Set("quantity", fmt.Sprint(rounded))
	v.Set("price", fmt.Sprint(px))
	resp, err := a.Post("/api/v3/order", v, true)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var out OrderResp
	if err := json.Unmarshal(bs, &out); err != nil {
		a.log.Infof("unmarshal OrderResp err:%v payload:%v", err, strings.TrimSpace(string(bs)))
		return "", err
	}

	if out.Code != 0 {
		return "", fmt.Errorf("limit order error: %v", out.Msg)
	}

	a.log.Info("<", strings.TrimSpace(string(bs)))
	return strconv.Itoa(out.OrderId), nil
}

func (a Api) MarketOrder(sym string, side data.Side, quoteQty *float64, baseQty *float64, _reduceOnly bool) (string, error) {
	exch, err := a.ExchangeInfo(sym)
	if err != nil {
		return "", err
	}

	var v = url.Values{}
	v.Set("side", strings.ToUpper(string(side)))
	v.Set("symbol", sym)
	v.Set("type", "MARKET")
	if quoteQty != nil {
		rounded := exch.RoundTickSize(sym, *quoteQty)
		if rounded < exch.MinNotional(sym) {
			a.log.Info("rounded quoteOrderQty smaller than MIN_NOTIONAL, skip place order")
			return "", nil
		}
		v.Set("quoteOrderQty", fmt.Sprint(rounded))
	}
	if baseQty != nil {
		rounded := exch.RoundLotSize(sym, math.Abs(*baseQty))
		if rounded == 0 {
			a.log.Info("rounded quantity smaller than LOT_SIZE, skip place order")
			return "", nil
		}
		tk, err := a.OrderBookTicker(sym)
		if err != nil {
			return "", err
		}
		if rounded*tk.BidPrice < exch.MinNotional(sym) {
			a.log.Info("rounded quoteOrderQty smaller than MIN_NOTIONAL, skip place order")
			return "", nil
		}
		v.Set("quantity", fmt.Sprint(rounded))
	}
	resp, err := a.Post("/api/v3/order", v, true)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var out OrderResp
	if err := json.Unmarshal(bs, &out); err != nil {
		return "", err
	}

	if out.Code != 0 {
		return "", fmt.Errorf("market order error: %v", out.Msg)
	}

	a.log.Info("<", strings.TrimSpace(string(bs)))
	return strconv.Itoa(out.OrderId), nil
}

func (a Api) GetOrder(sym, oid string) (*data.OrderStatus, error) {
	od, err := a.OrderStatus(sym, oid)
	if err != nil {
		return nil, err
	} else if od.Code != 0 {
		err = fmt.Errorf("fetch order status error code:%v msg:%v", od.Code, od.Msg)
		return nil, err
	}

	return &data.OrderStatus{
		Id:            strconv.Itoa(od.OrderId),
		Pair:          a.GetTradingPair(sym),
		Type:          data.OrderType(strings.ToLower(od.Type)),
		Side:          data.Side(strings.ToLower(od.Side)),
		Price:         od.Price,
		FilledSize:    od.ExecutedQty,
		RemainingSize: od.OrigQty - od.ExecutedQty,
		CreatedAt:     time.Unix(od.Time/1000, 0),
	}, nil
}

func (a Api) CancelOrder(sym, oid string) error {
	var v = url.Values{}
	v.Set("symbol", sym)
	v.Set("orderId", oid)

	resp, err := a.Delete("/api/v3/order", v, true)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var out OrderStatusResp
	if err := json.Unmarshal(bs, &out); err != nil {
		return err
	}

	if out.Code != 0 {
		return fmt.Errorf("%v", string(bs))
	}

	return nil
}
