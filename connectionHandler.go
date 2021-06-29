package main

import (
	"context"
	"fmt"
	"github.com/binance-exchange/go-binance"
	"github.com/go-kit/kit/log"
	"math"
	"os"
	"strconv"
	"time"
)

func createConnection(apiKey string, apiSecret string) binance.Binance {
	var logger log.Logger

	if __DEBUG__ {
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}

	hmacSigner := &binance.HmacSigner{
		Key: []byte(apiSecret),
	}
	ctx, _ := context.WithCancel(context.Background())
	binanceService := binance.NewAPIService(
		"https://www.binance.com",
		apiKey,
		hmacSigner,
		logger,
		ctx,
	)
	b := binance.NewBinance(binanceService)
	if __DEBUG__ {
		println(" => Finished Binance Authorization")
	}
	return b
}

func (pumpData *pumpData) createSellWaves(b binance.Binance) {
	currentPrice := getCurrentPrice(pumpData.asset, b)
	currentAssets := getAssetSize(pumpData.asset, b)
	if currentPrice == -1 || currentAssets == -1 {
		fmt.Println("error at sell wave creation")
		if __DEBUG__ {
			fmt.Println("current price: " + Float64ToString(currentPrice))
			fmt.Println("current asset size: " + Float64ToString(currentAssets))
		}
	}

	for i := 0; i < len(pumpData.sellLimits); i++ {

		sellPrice := pumpData.sellLimits[i][0]

		if sellPrice < currentPrice {
			sellPrice = currentPrice
		}

		sell, err := b.NewOrder(binance.NewOrderRequest{
			Symbol:    pumpData.asset,
			Side:      "SELL",
			Type:      "LIMIT",
			Quantity:  currentAssets * pumpData.sellLimits[i][1] / 100,
			Price:     sellPrice,
			Timestamp: time.Now(),
		})

		if err != nil {
			fmt.Println("error at sell limit creation!")
			if __DEBUG__ {
				fmt.Println("asset: " + pumpData.asset)
			}
		}

		fmt.Println("created sell-order for " + pumpData.asset + "! Order-ID: " + strconv.FormatInt(sell.OrderID, 10))
	}
}

func (pumpData *pumpData) pumpBuy(b binance.Binance, usdAmount string) bool {
	if !coinExists(pumpData.asset, b) {
		fmt.Println("error: asset does not exist - " + pumpData.asset)
		return false
	}

	buyAmount := math.Round(StringToFloat64(usdAmount) / getCurrentPrice(pumpData.asset, b))
	buyAmount = buyAmount - (buyAmount / 100)

	buy, err := b.NewOrder(binance.NewOrderRequest{
		Symbol:    pumpData.asset,
		Side:      "BUY",
		Type:      "LIMIT",
		Quantity:  buyAmount,
		Price:     pumpData.buyPrice,
		Timestamp: time.Now(),
	})

	if err != nil {
		fmt.Println("there was an unexpected error at pumpBuy")
		if __DEBUG__ {
			fmt.Println("asset: " + pumpData.asset + "\nUSD Amount: " + usdAmount)
			fmt.Println(err.Error())
		}
		return false
	}

	pumpData.transactionID = buy.ClientOrderID
	return true
}

func coinExists(asset string, b binance.Binance) bool {
	return getCurrentPrice(asset, b) != float64(-1)
}

func getCurrentPrice(asset string, b binance.Binance) float64 {
	if __DEBUG__ {
		println("Price request: " + asset)
	}

	price, requestError := b.Ticker24(binance.TickerRequest{Symbol: asset})
	if requestError != nil {
		if __DEBUG__ {
			fmt.Println("error at coin price request")
		}
		return -1
	}
	return price.BidPrice
}

func getAssetSize(asset string, bnc binance.Binance) float64 {
	balances, requestErr := bnc.Account(binance.AccountRequest{
		RecvWindow: 5 * time.Second,
		Timestamp:  time.Now(),
	})

	if requestErr != nil {
		fmt.Println("asset size request error!")
		if __DEBUG__ {
			fmt.Println("asset size request error! Coin: " + asset)
			fmt.Println("Error: " + requestErr.Error())
		}
		return -1
	}

	for _, requestedAsset := range balances.Balances {
		if requestedAsset.Asset == asset {
			return requestedAsset.Free
		}
	}
	return -1
}
