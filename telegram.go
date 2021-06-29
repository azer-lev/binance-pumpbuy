package main

import (
	"fmt"
	"github.com/binance-exchange/go-binance"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

func telegramConnection(b binance.Binance, apiKey string) {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		fmt.Println("telegram connection error")
		if __DEBUG__ {
			fmt.Println("DEBUG: APIKey => " + apiKey)
		}
		return
	}
	fmt.Println("Connection established: " + bot.Self.UserName)

	update := tgbotapi.NewUpdate(0)
	update.Timeout = 60

	updates, _ := bot.GetUpdatesChan(update)

	for currentUpdate := range updates {
		if currentUpdate.Message == nil || !strings.HasPrefix(currentUpdate.Message.Text, "/") {
			continue
		}

		if strings.HasPrefix(currentUpdate.Message.Text, "/pumpbuy") {
			// /pumpbuy ASSET 50 30000 30100,10,30200,90 29000
			if !strings.Contains(currentUpdate.Message.Text, " ") || len(strings.Split(currentUpdate.Message.Text, " ")) != __PUMPBUYARGUMENTS__ {
				bot.Send(tgbotapi.NewMessage(currentUpdate.Message.Chat.ID, "Unknown command usage: /pumpbuy => Use /help for more information!"))
				continue
			}
			messageParts := strings.Split(currentUpdate.Message.Text, " ")
			buyPriceConvert := StringToFloat64(messageParts[3])

			if buyPriceConvert == -1337 {
				bot.Send(tgbotapi.NewMessage(currentUpdate.Message.Chat.ID, "Unknown command usage: /pumpbuy => Use /help for more information!"))
				continue
			}

			pumpObj := pumpData{
				asset:          messageParts[1],
				buyPrice:       buyPriceConvert,
				sellLimits:     checkSellLimits(convertSellLimits(messageParts[4])),
				stopLoss:       StringToFloat64(messageParts[5]),
				percentageSell: false,
				transactionID:  "unset",
			}

			currentPump = pumpObj
			fmt.Println(currentPump.asset)
			fmt.Println(currentPump.buyPrice)
			fmt.Println(currentPump.stopLoss)

			currentPump.pumpBuy(b, messageParts[2])
			currentPump.createSellWaves(b)
			if __DEBUG__ {
				for _, row := range currentPump.sellLimits {
					for _, val := range row {
						fmt.Println(val)
					}
					fmt.Println()
				}
			}
		}
	}
}
