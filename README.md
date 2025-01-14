[![golangci-lint](https://github.com/oneart-dev/bybit/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/oneart-dev/bybit/actions/workflows/golangci-lint.yml)
[![test](https://github.com/oneart-dev/bybit/actions/workflows/test.yml/badge.svg)](https://github.com/oneart-dev/bybit/actions/workflows/test.yml)

# bybit

bybit is an bybit client for the Go programing language.

## Usage

### REST API

```
import "github.com/oneart-dev/bybit"

client := bybit.NewClient().WithAuth("your api key", "your api secret")
res, err := client.Future().InversePerpetual().Balance(bybit.CoinBTC)
// do as you want
```

### WebSocket API

for single use
```
import "github.com/oneart-dev/bybit"

wsClient := bybit.NewWebsocketClient()
svc, err := wsClient.Spot().V1().PublicV1()
if err != nil {
	return err
}
_, err = svc.SubscribeTrade(bybit.SymbolSpotBTCUSDT, func(response bybit.SpotWebsocketV1PublicV1TradeResponse) error {
	// do as you want
})
if err != nil {
	return err
}
svc.Start(context.Background())
```

for multiple use
```
import "github.com/oneart-dev/bybit"

wsClient := bybit.NewWebsocketClient()

executors := []bybit.WebsocketExecutor{}

svcRoot := wsClient.Spot().V1()
{
	svc, err := svcRoot.PublicV1()
	if err != nil {
		return err
	}
	_, err = svc.SubscribeTrade(bybit.SymbolSpotBTCUSDT, func(response bybit.SpotWebsocketV1PublicV1TradeResponse) error {
		// do as you want
	})
	if err != nil {
		return err
	}
	executors = append(executors, svc)
}
{
	svc, err := svcRoot.PublicV2()
	if err != nil {
		return err
	}
	_, err = svc.SubscribeTrade(bybit.SymbolSpotBTCUSDT, func(response bybit.SpotWebsocketV1PublicV2TradeResponse) error {
		// do as you want
	})
	if err != nil {
		return err
	}
	executors = append(executors, svc)
}

wsClient.Start(context.Background(), executors)
```

## Implemented

The following API endpoints have been implemented

### REST API

#### [Inverse Perpetual](https://bybit-exchange.github.io/docs/inverse)

##### Market Data Endpoints

- `/v2/public/orderBook/L2` Order Book
- `/v2/public/kline/list` Query Kline
- `/v2/public/tickers` Latest Information for Symbol
- `/v2/public/trading-records` Public Trading Records
- `/v2/public/symbols` Query Symbol
- `/v2/public/mark-price-kline` Query Mark Price Kline
- `/v2/public/index-price-kline` Query Index Price Kline
- `/v2/public/premium-index-kline` Query Premium Index Kline
- `/v2/public/open-interest` Open Interest
- `/v2/public/big-deal` Latest Big Deal
- `/v2/public/account-ratio` Long-Short Ratio

##### Account Data Endpoints

- `/v2/private/order/create` Place Active Order
- `/v2/private/order/list` Get Active Order
- `/v2/private/order/cancel` Cancel Active Order
- `/v2/private/position/list` My Position
- `/v2/private/position/leverage/save` Set Leverage

##### Wallet Data Endpoints

- `/v2/private/wallet/balance` Get Wallet Balance

#### [USDT Perpetual](https://bybit-exchange.github.io/docs/linear)

##### Market Data Endpoints

- `/v2/public/orderBook/L2` Order Book
- `/v2/public/tickers` Latest Information for Symbol
- `/v2/public/symbols` Query Symbol
- `/v2/public/open-interest` Open Interest
- `/v2/public/big-deal` Latest Big Deal
- `/v2/public/account-ratio` Long-Short Ratio

##### Account Data Endpoints

- `/private/linear/order/create` Place Active Order
- `/private/linear/order/cancel` Cancel Active Order
- `/private/linear/order/cancel-all` Cancel All Active Orders
- `/private/linear/position/list` My Position
- `/private/linear/position/set-leverage` Set Leverage
- `/private/linear/trade/execution/list` User Trade Records

##### Wallet Data Endpoints

- `/v2/private/wallet/balance` Get Wallet Balance

#### [Inverse Future](https://bybit-exchange.github.io/docs/inverse_futures)

##### Market Data Endpoints

- `/v2/public/orderBook/L2` Order Book
- `/v2/public/kline/list` Query Kline
- `/v2/public/tickers` Latest Information for Symbol
- `/v2/public/trading-records` Public Trading Records
- `/v2/public/symbols` Query Symbol
- `/v2/public/mark-price-kline` Query Index Price Kline
- `/v2/public/index-price-kline` Query Index Price Kline
- `/v2/public/open-interest` Open Interest
- `/v2/public/big-deal` Latest Big Deal
- `/v2/public/account-ratio` Long-Short Ratio

##### Wallet Data Endpoints

- `/v2/private/wallet/balance` Get Wallet Balance

#### [Spot](https://bybit-exchange.github.io/docs/spot)

##### Market Data Endpoints

- `/spot/v1/symbols` Query Symbol
- `/spot/quote/v1/depth` Order Book
- `/spot/quote/v1/depth/merged` Merged Order Book
- `/spot/quote/v1/trades` Public Trading Records
- `/spot/quote/v1/kline` Query Kline
- `/spot/quote/v1/ticker/24hr` Latest Information for Symbol
- `/spot/quote/v1/ticker/price` Last Traded Price
- `/spot/quote/v1/ticker/book_ticker` Best Bid/Ask Price

##### Account Data Endpoints

- `/spot/v1/order`
  - Place Active Order
  - Get Active Order
  - Cancel Active Order
  - Fast Cancel Active Order
- `/spot/v1/order/fast` Fast Cancel Active Order
- `/spot/order/batch-cancel` Batch Cancel Active Order
- `/spot/order/batch-fast-cancel` Batch Fast Cancel Active Order
- `/spot/order/batch-cancel-by-ids` Batch Cancel Active Order By IDs

### WebSocket API

#### [Spot v1](https://bybit-exchange.github.io/docs/spot/v1/#t-websocket)

##### Public Topics

- trade

##### Public Topics V2

- trade

##### Private Topics

- outboundAccountInfo
