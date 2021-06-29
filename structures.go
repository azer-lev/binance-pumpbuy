package main

import (
	"fmt"
	"strconv"
	"strings"
)

const __DEBUG__ = true
const __DEFAULTPATH__ = "config.lev"
const __CONFIGARGUMENTS__ = 3
const __PUMPBUYARGUMENTS__ = 6

var currentPump pumpData

type configData struct {
	apiKey      string
	apiSecret   string
	telegramKey string
}

type pumpData struct {
	asset          string
	buyPrice       float64
	sellLimits     [10][2]float64
	stopLoss       float64
	percentageSell bool

	transactionID string
}

func Float64ToString(inputNum float64) string {
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}

func StringToFloat64(input string) float64 {
	data, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return -1337
	}
	return data
}

func convertSellLimits(str string) [10][2]float64 {
	var res [10][2]float64
	split := strings.Split(str, ",")
	counter := 0
	for i := 0; i < len(split); i++ {
		if (i+1)%2 == 0 {
			res[counter][1] = StringToFloat64(split[i])
			counter++
		} else {
			res[counter][0] = StringToFloat64(split[i])
		}
	}
	return res
}

func checkSellLimits(uncheckedData [10][2]float64) [10][2]float64 {
	limitAddition := 0

	var fixedValues [10][2]float64

	for i := 0; i < len(uncheckedData); i++ {
		if uncheckedData[i][0] == 0 {
			if __DEBUG__ {
				fmt.Println("end of data")
			}
			continue
		}

		if limitAddition+int(uncheckedData[i][1]) > 100 {
			if limitAddition < 100 {
				lastCounter := len(fixedValues) - 1
				for fixedValues[lastCounter][0] == 0 && lastCounter > 0 {
					lastCounter--
				}
				fixedValues[lastCounter][1] += 100 - float64(limitAddition)
			}
			return fixedValues
		}

		fixedValues[i][0] = uncheckedData[i][0]
		fixedValues[i][1] = uncheckedData[i][1]
		limitAddition += int(uncheckedData[i][1])
	}
	return fixedValues
}
