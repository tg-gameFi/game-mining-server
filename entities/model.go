package entities

import "encoding/json"

type BNPrice struct {
	Symbol string      `json:"symbol"` // BTCUSDT
	Price  json.Number `json:"price"`
}

type CoinPrice struct {
	CoinSymbol string  `json:"coinSymbol"` // BTC
	Price      float64 `json:"price"`
}

type UserSession struct {
	Uid        int64
	IssuedAt   int64
	ExpiresSec int
}
