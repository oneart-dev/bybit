package bybit

import (
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/shopspring/decimal"
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
	ContractExecutionLists []ContractExecutionList `json:"list"`
}

type ContractExecutionList struct {
	OrderID          string          `json:"orderId"`
	OrderLinkID      string          `json:"orderLinkId"`
	Side             Side            `json:"side"`
	Symbol           string          `json:"symbol"`
	OrderPrice       decimal.Decimal `json:"orderPrice"`
	OrderQty         decimal.Decimal `json:"orderQty"`
	OrderType        OrderType       `json:"orderType"`
	FeeRate          decimal.Decimal `json:"feeRate"`
	ExecID           string          `json:"execId"`
	ExecPrice        decimal.Decimal `json:"execPrice"`
	ExecType         ExecType        `json:"execType"`
	ExecQty          decimal.Decimal `json:"execQty"`
	ExecFee          decimal.Decimal `json:"execFee"`
	ExecValue        decimal.Decimal `json:"execValue"`
	LeavesQty        decimal.Decimal `json:"leavesQty"`
	ClosedSize       decimal.Decimal `json:"closedSize"`
	LastLiquidityInd string          `json:"lastLiquidityInd"`
	TradeTimeMs      string          `json:"execTime"`
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
	Balance []ContractBalance `json:"list"`
}

// Balance :
type ContractBalance struct {
	Coin             string          `json:"coin"`
	Equity           decimal.Decimal `json:"equity"`
	AvailableBalance decimal.Decimal `json:"availableBalance"`
	UsedMargin       decimal.Decimal `json:"used_margin"`
	OrderMargin      decimal.Decimal `json:"orderMargin"`
	PositionMargin   decimal.Decimal `json:"positionMargin"`
	OccClosingFee    decimal.Decimal `json:"occClosingFee"`
	OccFundingFee    decimal.Decimal `json:"occFundingFee"`
	WalletBalance    decimal.Decimal `json:"walletBalance"`
	RealisedPnl      decimal.Decimal `json:"realised_pnl"`
	UnrealisedPnl    decimal.Decimal `json:"unrealisedPnl"`
	CumRealisedPnl   decimal.Decimal `json:"cumRealisedPnl"`
	GivenCash        decimal.Decimal `json:"givenCash"`
	ServiceCash      decimal.Decimal `json:"serviceCash"`
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
	MaxTradingQty string `json:"maxTradingQty"`
	MinTradingQty string `json:"minTradingQty"`
	QtyStep       string `json:"qtyStep"`
}

// LeverageFilter :
type ContractLeverageFilter struct {
	MinLeverage  string `json:"minLeverage"`
	MaxLeverage  string `json:"maxLeverage"`
	LeverageStep string `json:"leverageStep"`
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
	Symbol            string        `json:"symbol"`
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
	OpenInterest      string        `json:"openInterest"`
	OpenValue         string        `json:"open_value"`
	TotalTurnover     string        `json:"turnover24h"`
	Turnover24h       string        `json:"turnover_24h"`
	TotalVolume       string        `json:"total_volume"`
	Volume24h         string        `json:"volume24h"`
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
