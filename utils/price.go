package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"game-mining-server/entities"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const retryTimes = 3
const baseUrl = "https://api.binance.com"
const quoteSymbol = "USDT"

var httpClient = &http.Client{Timeout: 10 * time.Second}

// GetCoinUSDPrice Try to fetch coin price from third party service
func GetCoinUSDPrice(coinSymbols []string) ([]*entities.CoinPrice, error) {
	var pairSymbols []string
	for _, coinSymbol := range coinSymbols {
		if coinSymbol == quoteSymbol {
			continue
		} else {
			pairSymbols = append(pairSymbols, coinSymbol+quoteSymbol)
		}
	}

	params := url.Values{}
	params.Set("symbols", fmt.Sprintf("[\"%s\"]", strings.Join(pairSymbols, "\",\"")))

	priceRes, e0 := bnRequestGet("/api/v3/ticker/price", &params)
	if e0 != nil {
		return nil, e0
	}
	var bnPrices []*entities.BNPrice
	if e1 := json.Unmarshal(priceRes, &bnPrices); e1 != nil {
		return nil, e1
	}

	var coinPrices []*entities.CoinPrice
	for _, priceItem := range bnPrices {
		symbolSplits := strings.Split(priceItem.Symbol, quoteSymbol)
		if len(symbolSplits) != 2 {
			continue
		}
		price, _ := priceItem.Price.Float64()
		coinPrices = append(coinPrices, &entities.CoinPrice{
			CoinSymbol: symbolSplits[0],
			Price:      price,
		})
	}
	return coinPrices, nil
}

func bnRequestGet(path string, params *url.Values) ([]byte, error) {
	finalUrl, e0 := url.Parse(baseUrl + path)
	if e0 != nil {
		return nil, e0
	}
	finalUrl.RawQuery = params.Encode()
	for i := 0; i < retryTimes; i++ {
		if res, e1 := httpClient.Get(finalUrl.String()); e1 != nil {
			var err net.Error
			if errors.As(e1, &err) && err.Timeout() {
				time.Sleep(500 * time.Millisecond)
			} else {
				return nil, e1 // unknown error, throw it
			}
		} else {
			body, _ := io.ReadAll(res.Body)
			_ = res.Body.Close()
			if strings.Contains(string(body), "Timestamp for this request is outside of the recvWindow") {
				time.Sleep(500 * time.Millisecond)
			} else {
				return body, nil
			}
		}
	}
	return nil, fmt.Errorf("price request GET %s failed after %d retries", path, retryTimes)
}
