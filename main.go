package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("there was an unexpected error at pumpBuy")
	if !configReady(__DEFAULTPATH__) {
		if !pathExists(__DEFAULTPATH__) {
			createDefault(__DEFAULTPATH__)
		}
		fmt.Println("Please make sure to update the " + __DEFAULTPATH__)
		os.Exit(100)
	}
	cfg := getData(__DEFAULTPATH__)
	bnc := createConnection(cfg.apiKey, cfg.apiSecret)
	getAssetSize("BTC", bnc)
	telegramConnection(bnc, cfg.telegramKey)
}
