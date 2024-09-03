package middleware

import (
	"game-mining-server/app"
	"game-mining-server/caches"
	"game-mining-server/entities"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"net/http"
)

func LimitIp30PerMinMiddleware() gin.HandlerFunc {
	return createIpRateLimiter(30)
}

func LimitIp60PerMinMiddleware() gin.HandlerFunc {
	return createIpRateLimiter(60)
}

func LimitIp120PerMinMiddleware() gin.HandlerFunc {
	return createIpRateLimiter(120)
}

func LimitIp240PerMinMiddleware() gin.HandlerFunc { return createIpRateLimiter(240) }

func LimitIp480PerMinMiddleware() gin.HandlerFunc { return createIpRateLimiter(480) }

func createIpRateLimiter(rate int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := app.Cache().RateLimiter.Allow(ctx, caches.GenRateLimitCacheKey("IP", ctx.RemoteIP()), redis_rate.PerMinute(rate))
		if err != nil {
			ctx.Abort()
			return
		}
		if res.Allowed == 0 {
			ctx.JSON(http.StatusTooManyRequests, entities.ResFailed(entities.ErrTooManyRequests, "too many requests"))
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}
