package bybit

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-querystring/query"
)

// V5OrderServiceI :
type V5OrderServiceI interface {
	CreateOrder(V5CreateOrderParam) (*V5CreateOrderResponse, error)
	CancelOrder(V5CancelOrderParam) (*V5CancelOrderResponse, error)
	GetOpenOrders(V5GetOpenOrdersParam) (*V5GetOpenOrdersResponse, error)
	GetExecutionList(V5GetExecutionListParam) (*V5GetExecutionListResponse, error)
	GetOrderList(param V5GetOrderListParam) (*V5GetOrderListResponse, error)
	GetClosedPnl(param V5GetClosedPnlParam) (*V5GetClosedPnlResponse, error)
}

// V5OrderService :
type V5OrderService struct {
	client *Client
}

// V5CreateOrderParam :
type V5CreateOrderParam struct {
	Category  CategoryV5 `json:"category"`
	Symbol    SymbolV5   `json:"symbol"`
	Side      Side       `json:"side"`
	OrderType OrderType  `json:"orderType"`
	Qty       string     `json:"qty"`

	IsLeverage            *IsLeverage       `json:"isLeverage,omitempty"`
	Price                 *string           `json:"price,omitempty"`
	TriggerDirection      *TriggerDirection `json:"triggerDirection,omitempty"`
	OrderFilter           *OrderFilter      `json:"orderFilter,omitempty"` // If not passed, Order by default
	TriggerPrice          *string           `json:"triggerPrice,omitempty"`
	TriggerBy             *TriggerBy        `json:"triggerBy,omitempty"`
	OrderIv               *string           `json:"orderIv,omitempty"`     // option only.
	TimeInForce           *TimeInForce      `json:"timeInForce,omitempty"` // If not passed, GTC is used by default
	PositionIdx           *PositionIdx      `json:"positionIdx,omitempty"` // Under hedge-mode, this param is required
	OrderLinkID           *string           `json:"orderLinkId,omitempty"`
	TakeProfit            *string           `json:"takeProfit,omitempty"`
	StopLoss              *string           `json:"stopLoss,omitempty"`
	TpTriggerBy           *TriggerBy        `json:"tpTriggerBy,omitempty"`
	SlTriggerBy           *TriggerBy        `json:"slTriggerBy,omitempty"`
	ReduceOnly            *bool             `json:"reduce_only,omitempty"`
	CloseOnTrigger        *bool             `json:"closeOnTrigger,omitempty"`
	MarketMakerProtection *bool             `json:"mmp,omitempty"` // option only
}

// V5CreateOrderResponse :
type V5CreateOrderResponse struct {
	CommonV5Response `json:",inline"`
	Result           V5CreateOrderResult `json:"result"`
}

// V5CreateOrderResult :
type V5CreateOrderResult struct {
	OrderID     string `json:"orderId"`
	OrderLinkID string `json:"orderLinkId"`
}

// CreateOrder :
func (s *V5OrderService) CreateOrder(param V5CreateOrderParam) (*V5CreateOrderResponse, error) {
	var res V5CreateOrderResponse

	body, err := json.Marshal(param)
	if err != nil {
		return &res, fmt.Errorf("json marshal: %w", err)
	}

	if err := s.client.postV5JSON("/v5/order/create", body, &res); err != nil {
		return &res, err
	}

	return &res, nil
}

// V5CancelOrderParam :
type V5CancelOrderParam struct {
	Category CategoryV5 `json:"category"`
	Symbol   SymbolV5   `json:"symbol"`

	OrderID     *string      `json:"orderId,omitempty"`
	OrderLinkID *string      `json:"orderLinkId,omitempty"`
	OrderFilter *OrderFilter `json:"orderFilter,omitempty"` // If not passed, Order by default
}

// V5CancelOrderResponse :
type V5CancelOrderResponse struct {
	CommonV5Response `json:",inline"`
	Result           V5CancelOrderResult `json:"result"`
}

// V5CancelOrderResult :
type V5CancelOrderResult struct {
	OrderID     string `json:"orderId"`
	OrderLinkID string `json:"orderLinkId"`
}

// CancelOrder :
func (s *V5OrderService) CancelOrder(param V5CancelOrderParam) (*V5CancelOrderResponse, error) {
	var res V5CancelOrderResponse

	if param.OrderID == nil && param.OrderLinkID == nil {
		return nil, fmt.Errorf("either OrderID or OrderLinkID needed")
	}

	body, err := json.Marshal(param)
	if err != nil {
		return &res, fmt.Errorf("json marshal: %w", err)
	}

	if err := s.client.postV5JSON("/v5/order/cancel", body, &res); err != nil {
		return &res, err
	}

	return &res, nil
}

// V5GetOpenOrdersParam :
type V5GetOpenOrdersParam struct {
	Category CategoryV5 `url:"category"`

	Symbol      *SymbolV5    `url:"symbol,omitempty"`
	BaseCoin    *Coin        `url:"baseCoin,omitempty"`
	SettleCoin  *Coin        `url:"settleCoin,omitempty"`
	OrderID     *string      `url:"orderId,omitempty"`
	OrderLinkID *string      `url:"orderLinkId,omitempty"`
	OpenOnly    *int         `url:"openOnly,omitempty"`
	OrderFilter *OrderFilter `url:"orderFilter,omitempty"` // If not passed, Order by default
	Limit       *int         `url:"limit,omitempty"`
	Cursor      *string      `url:"cursor,omitempty"`
}

// V5GetOpenOrdersResponse :
type V5GetOpenOrdersResponse struct {
	CommonV5Response `json:",inline"`
	Result           V5GetOpenOrdersResult `json:"result"`
}

// V5GetOpenOrdersResult :
type V5GetOpenOrdersResult struct {
	Category       CategoryV5       `json:"category"`
	NextPageCursor string           `json:"nextPageCursor"`
	List           []V5GetOpenOrder `json:"list"`
}

type V5GetOpenOrder struct {
	Symbol             SymbolV5    `json:"symbol"`
	OrderType          OrderType   `json:"orderType"`
	OrderLinkID        string      `json:"orderLinkId"`
	OrderID            string      `json:"orderId"`
	CancelType         string      `json:"cancelType"`
	AvgPrice           string      `json:"avgPrice"`
	StopOrderType      string      `json:"stopOrderType"`
	LastPriceOnCreated string      `json:"lastPriceOnCreated"`
	OrderStatus        OrderStatus `json:"orderStatus"`
	TakeProfit         string      `json:"takeProfit"`
	CumExecValue       string      `json:"cumExecValue"`
	TriggerDirection   int         `json:"triggerDirection"`
	IsLeverage         string      `json:"isLeverage"`
	RejectReason       string      `json:"rejectReason"`
	Price              string      `json:"price"`
	OrderIv            string      `json:"orderIv"`
	CreatedTime        string      `json:"createdTime"`
	TpTriggerBy        string      `json:"tpTriggerBy"`
	PositionIdx        int         `json:"positionIdx"`
	TimeInForce        TimeInForce `json:"timeInForce"`
	LeavesValue        string      `json:"leavesValue"`
	UpdatedTime        string      `json:"updatedTime"`
	Side               Side        `json:"side"`
	TriggerPrice       string      `json:"triggerPrice"`
	CumExecFee         string      `json:"cumExecFee"`
	LeavesQty          string      `json:"leavesQty"`
	SlTriggerBy        string      `json:"slTriggerBy"`
	CloseOnTrigger     bool        `json:"closeOnTrigger"`
	CumExecQty         string      `json:"cumExecQty"`
	ReduceOnly         bool        `json:"reduceOnly"`
	Qty                string      `json:"qty"`
	StopLoss           string      `json:"stopLoss"`
	TriggerBy          TriggerBy   `json:"triggerBy"`
}

// GetOpenOrders :
func (s *V5OrderService) GetOpenOrders(param V5GetOpenOrdersParam) (*V5GetOpenOrdersResponse, error) {
	var res V5GetOpenOrdersResponse

	if param.Category == "" {
		return nil, fmt.Errorf("Category needed")
	}

	queryString, err := query.Values(param)
	if err != nil {
		return nil, err
	}

	if err := s.client.getV5Privately("/v5/order/realtime", queryString, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type V5GetExecutionListParam struct {
	Category CategoryV5 `url:"category"`

	StartTime   *int      `url:"startTime,omitempty"`
	EndTime     *int      `url:"endTime,omitempty"`
	ExecType    *ExecType `url:"execType,omitempty"`
	Symbol      *SymbolV5 `url:"symbol,omitempty"`
	BaseCoin    *Coin     `url:"baseCoin,omitempty"`
	OrderID     *string   `url:"orderId,omitempty"`
	OrderLinkID *string   `url:"orderLinkId,omitempty"`
	Limit       *int      `url:"limit,omitempty"`
	Cursor      *string   `url:"cursor,omitempty"`
}

// V5GetOpenOrdersResponse :
type V5GetExecutionListResponse struct {
	CommonV5Response `json:",inline"`
	Result           V5GetExecutionListResult `json:"result"`
}

// V5GetOpenOrdersResult :
type V5GetExecutionListResult struct {
	Category       CategoryV5            `json:"category"`
	NextPageCursor string                `json:"nextPageCursor"`
	List           []V5GetExecutionOrder `json:"list"`
}

type V5GetExecutionOrder struct {
	Symbol        SymbolV5  `json:"symbol"`
	OrderType     OrderType `json:"orderType"`
	OrderLinkID   string    `json:"orderLinkId"`
	Side          Side      `json:"side"`
	OrderID       string    `json:"orderId"`
	StopOrderType string    `json:"stopOrderType"`
	LeavesQty     string    `json:"leavesQty"`
	ExecTime      string    `json:"execTime"`
	IsMaker       bool      `json:"isMaker"`
	ExecFee       string    `json:"execFee"`
	FeeRate       string    `json:"feeRate"`
	ExecID        string    `json:"execId"`
	MarkPrice     string    `json:"markPrice"`
	ExecPrice     string    `json:"execPrice"`
	OrderQty      string    `json:"orderQty"`
	OrderPrice    string    `json:"orderPrice"`
	ExecValue     string    `json:"execValue"`
	ExecType      ExecType  `json:"execType"`
	ExecQty       string    `json:"execQty"`
}

// GetExecutionList :
func (s *V5OrderService) GetExecutionList(param V5GetExecutionListParam) (*V5GetExecutionListResponse, error) {
	var res V5GetExecutionListResponse

	if param.Category == "" {
		return nil, fmt.Errorf("category needed")
	}

	queryString, err := query.Values(param)
	if err != nil {
		return nil, err
	}

	if err := s.client.getV5Privately("/v5/execution/list", queryString, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type V5GetOrderListParam struct {
	Category CategoryV5 `url:"category"`

	StartTime   *int      `url:"startTime,omitempty"`
	EndTime     *int      `url:"endTime,omitempty"`
	ExecType    *ExecType `url:"execType,omitempty"`
	Symbol      *SymbolV5 `url:"symbol,omitempty"`
	BaseCoin    *Coin     `url:"baseCoin,omitempty"`
	OrderID     *string   `url:"orderId,omitempty"`
	OrderLinkID *string   `url:"orderLinkId,omitempty"`
	Limit       *int      `url:"limit,omitempty"`
	Cursor      *string   `url:"cursor,omitempty"`
}

// V5GetOpenOrdersResponse :
type V5GetOrderListResponse struct {
	CommonV5Response `json:",inline"`
	Result           V5GetOrderListResult `json:"result"`
}

// V5GetOpenOrdersResult :
type V5GetOrderListResult struct {
	Category       CategoryV5   `json:"category"`
	NextPageCursor string       `json:"nextPageCursor"`
	List           []V5GetOrder `json:"list"`
}

type V5GetOrder struct {
	Symbol      SymbolV5  `json:"symbol"`
	OrderType   OrderType `json:"orderType"`
	OrderLinkID string    `json:"orderLinkId"`
	OrderID     string    `json:"orderId"`

	AvgPrice           string `json:"avgPrice"`
	StopOrderType      string `json:"stopOrderType"`
	LastPriceOnCreated string `json:"lastPriceOnCreated"`
	OrderStatus        string `json:"orderStatus"`
	TakeProfit         string `json:"takeProfit"`
	CumExecValue       string `json:"cumExecValue"`
	TriggerDirection   int    `json:"triggerDirection"`
	BlockTradeID       string `json:"blockTradeId"`
	RejectReason       string `json:"rejectReason"`
	IsLeverage         string `json:"isLeverage"`
	Price              string `json:"price"`
	OrderIV            string `json:"orderIv"`
	CreatedTime        string `json:"createdTime"`
	TPTriggerBy        string `json:"tpTriggerBy"`
	PositionIdx        int    `json:"positionIdx"`
	TimeInForce        string `json:"timeInForce"`
	LeavesValue        string `json:"leavesValue"`
	UpdatedTime        string `json:"updatedTime"`
	Side               Side   `json:"side"`
	TriggerPrice       string `json:"triggerPrice"`
	CumExecFee         string `json:"cumExecFee"`
	SLTriggerBy        string `json:"slTriggerBy"`
	LeavesQty          string `json:"leavesQty"`
	CloseOnTrigger     bool   `json:"closeOnTrigger"`
	CumExecQty         string `json:"cumExecQty"`
	ReduceOnly         bool   `json:"reduceOnly"`
	Qty                string `json:"qty"`
	StopLoss           string `json:"stopLoss"`
	TriggerBy          string `json:"triggerBy"`
}

// GetOrderList :
func (s *V5OrderService) GetOrderList(param V5GetOrderListParam) (*V5GetOrderListResponse, error) {
	var res V5GetOrderListResponse

	if param.Category == "" {
		return nil, fmt.Errorf("category needed")
	}

	queryString, err := query.Values(param)
	if err != nil {
		return nil, err
	}

	if err := s.client.getV5Privately("/v5/order/history", queryString, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type V5GetClosedPnlParam struct {
	Category CategoryV5 `url:"category"`

	StartTime *int      `url:"startTime,omitempty"`
	EndTime   *int      `url:"endTime,omitempty"`
	Symbol    *SymbolV5 `url:"symbol,omitempty"`
	Limit     *int      `url:"limit,omitempty"`
	Cursor    *string   `url:"cursor,omitempty"`
}

// V5GetOpenOrdersResponse :
type V5GetClosedPnlResponse struct {
	CommonV5Response `json:",inline"`
	Result           V5GetClosedPnlResult `json:"result"`
}

// V5GetOpenOrdersResult :
type V5GetClosedPnlResult struct {
	Category       CategoryV5   `json:"category"`
	NextPageCursor string       `json:"nextPageCursor"`
	List           []V5GetOrder `json:"list"`
}

/*
{
                "symbol": "ETHPERP",
                "orderType": "Market",
                "leverage": "3",
                "updatedTime": "1672214887236",
                "side": "Sell",
                "orderId": "5a373bfe-188d-4913-9c81-d57ab5be8068",
                "closedPnl": "-47.4065323",
                "avgEntryPrice": "1194.97516667",
                "qty": "3",
                "cumEntryValue": "3584.9255",
                "createdTime": "1672214887231423699",
                "orderPrice": "1122.95",
                "closedSize": "3",
                "avgExitPrice": "1180.59833333",
                "execType": "Trade",
                "fillCount": "4",
                "cumExitValue": "3541.795"
            }
*/

type V5GetClosedPnl struct {
	Symbol    SymbolV5  `json:"symbol"`
	OrderType OrderType `json:"orderType"`
	OrderID   string    `json:"orderId"`

	Leverage    string `json:"leverage"`
	UpdatedTime string `json:"updatedTime"`
	Side        Side   `json:"side"`
	ClosedPnl   string `json:"closedPnl"`
	AvgEntryPrice string `json:"avgEntryPrice"`
	Qty         string `json:"qty"`
	CumEntryValue string `json:"cumEntryValue"`
	CreatedTime string `json:"createdTime"`
	OrderPrice  string `json:"orderPrice"`
	ClosedSize  string `json:"closedSize"`
	AvgExitPrice string `json:"avgExitPrice"`
	ExecType    string `json:"execType"`
	FillCount   string `json:"fillCount"`
	CumExitValue string `json:"cumExitValue"`
}

// GetOrderList :
func (s *V5OrderService) GetClosedPnl(param V5GetClosedPnlParam) (*V5GetClosedPnlResponse, error) {
	var res V5GetClosedPnlResponse

	if param.Category == "" {
		return nil, fmt.Errorf("category needed")
	}

	queryString, err := query.Values(param)
	if err != nil {
		return nil, err
	}

	if err := s.client.getV5Privately("/v5/position/closed-pnl", queryString, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
