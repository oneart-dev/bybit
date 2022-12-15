package bybit

import (
	"net/url"

	"github.com/google/go-querystring/query"
)

// FutureContractServiceI :
type FutureContractServiceI interface {
	// Market Data Endpoints
	Tickers(category string, symbol string) (*ContractTickersResponse, error)
	Symbols() (*ContractSymbolsResponse, error)
	ListKline(ContractListKlineParam) (*ContractListKlineResponse, error)

	// Account Data Endpoints
	ContractExecutionHistoryList(param ContractExecutionHistoryListParam) (*ContractExecutionHistoryListResponse, error)

	// Wallet Data Endpoints
	Balance(Coin) (*ContractBalanceResponse, error)
}

// FutureContractService :
type FutureContractService struct {
	client *Client

	*FutureCommonService
}

type ContractListKlineParam struct {
	Category string        `url:"category"`
	Symbol   SymbolInverse `url:"symbol"`
	Interval Interval      `url:"interval"`
	Start    int           `url:"start"`
	End      int           `url:"end"`

	Limit *int `url:"limit,omitempty"`
}

type ContractListKlineResponse struct {
	CommonResponse `json:",inline"`
	Result         ContractListKlineResult `json:"result"`
}

// ListKlineResult :
type ContractListKlineResult struct {
	Category string        `json:"category"`
	Symbol   SymbolInverse `json:"symbol"`
	Interval string        `json:"interval"`
	List     [][]string    `json:"list"`
}

// ListKline :
func (s *FutureContractService) ListKline(param ContractListKlineParam) (*ContractListKlineResponse, error) {
	var res ContractListKlineResponse

	queryString, err := query.Values(param)
	if err != nil {
		return nil, err
	}

	if err := s.client.getPublicly("/derivatives/v3/public/kline", queryString, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ContractExecutionHistoryListResponse :
type ContractExecutionHistoryListResponse struct {
	CommonResponse `json:",inline"`
	Result         ContractExecutionHistoryListResult `json:"result"`
}

// ContractExecutionHistoryListResult :
type ContractExecutionHistoryListResult struct {
	PageToken              string                  `json:"nextPageCursor"`
	ContractExecutionLists []ContractExecutionList `json:"data"`
}

type ContractExecutionList struct {
	OrderID          string     `json:"orderId"`
	OrderLinkID      string     `json:"orderLinkId"`
	Side             Side       `json:"side"`
	Symbol           SymbolUSDT `json:"symbol"`
	OrderPrice       float64    `json:"orderPrice"`
	OrderQty         float64    `json:"orderQty"`
	OrderType        OrderType  `json:"orderType"`
	FeeRate          float64    `json:"feeRate"`
	ExecID           string     `json:"execId"`
	ExecPrice        float64    `json:"execPrice"`
	ExecType         ExecType   `json:"execType"`
	ExecQty          float64    `json:"execQty"`
	ExecFee          float64    `json:"execFee"`
	ExecValue        float64    `json:"execValue"`
	LeavesQty        float64    `json:"leavesQty"`
	ClosedSize       float64    `json:"closedSize"`
	LastLiquidityInd string     `json:"lastLiquidityInd"`
	TradeTimeMs      int64      `json:"execTime"`
}

// ContractExecutionHistoryListParam :
type ContractExecutionHistoryListParam struct {
	Symbol SymbolUSDT `url:"symbol"`

	StartTime *int      `url:"start_time,omitempty"`
	EndTime   *int      `url:"end_time,omitempty"`
	ExecType  *ExecType `url:"exec_type,omitempty"`
	PageToken *string   `url:"cursor,omitempty"`
	Limit     *int      `url:"limit,omitempty"`
}

// ContractExecutionHistoryList :
func (s *FutureContractService) ContractExecutionHistoryList(param ContractExecutionHistoryListParam) (*ContractExecutionHistoryListResponse, error) {
	var res ContractExecutionHistoryListResponse

	queryString, err := query.Values(param)
	if err != nil {
		return nil, err
	}

	if err := s.client.getPrivately("/contract/v3/private/execution/list", queryString, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// BalanceResponse :
type ContractBalanceResponse struct {
	CommonResponse `json:",inline"`
	Result         ContractBalanceResult `json:"result"`
}

// BalanceResult :
type ContractBalanceResult struct {
	Balance []Balance `json:"list"`
}

// Balance :
type ContractBalance struct {
	Coin             string  `json:"coin"`
	Equity           float64 `json:"equity"`
	AvailableBalance float64 `json:"availableBalance"`
	UsedMargin       float64 `json:"used_margin"`
	OrderMargin      float64 `json:"orderMargin"`
	PositionMargin   float64 `json:"positionMargin"`
	OccClosingFee    float64 `json:"occClosingFee"`
	OccFundingFee    float64 `json:"occFundingFee"`
	WalletBalance    float64 `json:"walletBalance"`
	RealisedPnl      float64 `json:"realised_pnl"`
	UnrealisedPnl    float64 `json:"unrealisedPnl"`
	CumRealisedPnl   float64 `json:"cumRealisedPnl"`
	GivenCash        float64 `json:"givenCash"`
	ServiceCash      float64 `json:"serviceCash"`
}

// Balance :
func (s *FutureContractService) Balance(coin Coin) (*ContractBalanceResponse, error) {
	var res ContractBalanceResponse

	query := url.Values{}
	query.Add("coin", string(coin))
	if err := s.client.getPrivately("/contract/v3/private/account/wallet/balance", query, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// SymbolsResponse :
type ContractSymbolsResponse struct {
	CommonResponse `json:",inline"`
	Result         ContractSymbolData `json:"result"`
}

type ContractSymbolData struct {
	Category string           `json:"category"`
	List     []ContractSymbol `json:"list"`
	Page     string           `json:"nextPageCursor"`
}

// SymbolsResult :
type ContractSymbol struct {
	Name           string                 `json:"symbol"`
	Status         string                 `json:"status"`
	BaseCurrency   string                 `json:"baseCoin"`
	QuoteCurrency  string                 `json:"quoteCoin"`
	PriceScale     string                 `json:"priceScale"`
	TakerFee       string                 `json:"taker_fee"`
	MakerFee       string                 `json:"maker_fee"`
	LeverageFilter ContractLeverageFilter `json:"leverageFilter"`
	PriceFilter    ContractPriceFilter    `json:"priceFilter"`
	LotSizeFilter  ContractLotSizeFilter  `json:"lotSizeFilter"`
}

// PriceFilter :
type ContractPriceFilter struct {
	MinPrice string `json:"minPrice"`
	MaxPrice string `json:"maxPrice"`
	TickSize string `json:"tickSize"`
}

// LotSizeFilter :
type ContractLotSizeFilter struct {
	MaxTradingQty float64 `json:"maxTradingQty"`
	MinTradingQty float64 `json:"minTradingQty"`
	QtyStep       float64 `json:"qtyStep"`
}

// LeverageFilter :
type ContractLeverageFilter struct {
	MinLeverage  float64 `json:"minLeverage"`
	MaxLeverage  float64 `json:"maxLeverage"`
	LeverageStep string  `json:"leverageStep"`
}

// Symbols :
func (s *FutureContractService) Symbols() (*ContractSymbolsResponse, error) {
	var res ContractSymbolsResponse

	if err := s.client.getPublicly("/derivatives/v3/public/instruments-info", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// TickersResponse :
type ContractTickersResponse struct {
	CommonResponse `json:",inline"`
	Result         ContractTickersResult `json:"result"`
}

// TickersResult :
type ContractTickersResult struct {
	Category string           `json:"category"`
	List     []ContractTicker `json:"list"`
}

type ContractTicker struct {
	Symbol            SymbolInverse `json:"symbol"`
	BidPrice          string        `json:"bidPrice"`
	AskPrice          string        `json:"askPrice"`
	LastPrice         string        `json:"lastPrice"`
	LastTickDirection TickDirection `json:"lastTickDirection"`
	PrevPrice24h      string        `json:"prevPrice24h"`
	Price24hPcnt      string        `json:"price24hPcnt"`
	HighPrice24h      string        `json:"highPrice24h"`
	LowPrice24h       string        `json:"lowPrice24h"`
	PrevPrice1h       string        `json:"prevPrice1h"`
	MarkPrice         string        `json:"markPrice"`
	IndexPrice        string        `json:"indexPrice"`
	OpenInterest      float64       `json:"openInterest"`
	OpenValue         string        `json:"open_value"`
	TotalTurnover     string        `json:"turnover24h"`
	Turnover24h       string        `json:"turnover_24h"`
	TotalVolume       float64       `json:"total_volume"`
	Volume24h         float64       `json:"volume24h"`
	FundingRate       string        `json:"fundingRate"`
	NextFundingTime   string        `json:"nextFundingTime"`
}

// Tickers :
func (s *FutureContractService) Tickers(category string, symbol string) (*ContractTickersResponse, error) {
	var res ContractTickersResponse

	query := url.Values{}
	query.Add("symbol", symbol)
	query.Add("category", category)

	if err := s.client.getPublicly("/derivatives/v3/public/tickers", query, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
