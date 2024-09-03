package api

import (
	"game-mining-server/app"
	"game-mining-server/caches"
	"game-mining-server/entities"
	"game-mining-server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func GetCoinPrice(c *gin.Context) {
	var params entities.CoinPriceParam
	if e0 := c.ShouldBindQuery(&params); e0 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, e0.Error()))
		return
	}

	coinSymbols := strings.Split(params.CoinSymbols, ",")

	var coinPriceMap = make(map[string]float64)
	var needRequestCoinSymbols = make([]string, 0)
	for _, coinSymbol := range coinSymbols {
		cachedPrice, e1 := app.Cache().GetString(caches.GenCoinPriceCacheKey(params.FiatSymbol, coinSymbol))
		if e1 == nil && cachedPrice != "" { // cache hit, just return cache price
			priceFloat, _ := strconv.ParseFloat(cachedPrice, 64)
			coinPriceMap[coinSymbol] = priceFloat
		} else {
			needRequestCoinSymbols = append(needRequestCoinSymbols, coinSymbol)
		}
	}

	if len(needRequestCoinSymbols) > 0 {
		latestCoinPrices, e2 := utils.GetCoinUSDPrice(needRequestCoinSymbols)
		if e2 != nil {
			c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, e2.Error()))
			return
		}

		for _, price := range latestCoinPrices {
			_ = app.Cache().SetString(caches.GenCoinPriceCacheKey(params.FiatSymbol, price.CoinSymbol), strconv.FormatFloat(price.Price, 'E', -1, 64), 60)
			coinPriceMap[price.CoinSymbol] = price.Price
		}
	}

	c.JSON(http.StatusOK, entities.ResSuccess(coinPriceMap))
}
