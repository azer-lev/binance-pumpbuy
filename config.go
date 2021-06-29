package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func pathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func configReady(path string) bool {
	if pathExists(path) {
		cfg := getData(path)
		if cfg.apiKey == "InsertyourBinanceAPIKey!" || cfg.apiSecret == "InsertyourBinanceAPISecret!" {
			return false
		}
		return true
	}
	return false
}

//Creates a default file, if no file is existing
func createDefault(path string) {
	if !pathExists(path) {
		os.Create(path)
		defaultDate := "Binance API Key = Insert your Binance API Key!\n" +
			"Binance API Secret = Insert your Binance API Secret!\n" +
			"Telegram API Key = Insert your Telegram API Key!\n"

		fileWriter := ioutil.WriteFile(path, []byte(defaultDate), 0644)
		if fileWriter != nil {
			fmt.Println(fileWriter.Error())
		}
		return
	}
	if __DEBUG__ {
		fmt.Println("path does already exist! Can't create default file")
	}
}

func getData(path string) configData {
	if !pathExists(path) {
		if __DEBUG__ {
			fmt.Println("Path does not exist! Can't read data!")
		}
		return configData{}
	}

	data, err := ioutil.ReadFile(path)

	if err != nil {
		if __DEBUG__ {
			fmt.Println("Error at reading file-data!")
		}
		return configData{}
	}

	var lines = strings.Split(string(data), "\n")

	if len(lines) < __CONFIGARGUMENTS__ {
		if __DEBUG__ {
			fmt.Println("Config file-data corrupt!")
		}
		return configData{}
	}

	return configData{
		apiKey:      strings.Replace(strings.Split(lines[0], "=")[1], " ", "", -1),
		apiSecret:   strings.Replace(strings.Split(lines[1], "=")[1], " ", "", -1),
		telegramKey: strings.Replace(strings.Split(lines[2], "=")[1], " ", "", -1),
	}
}
