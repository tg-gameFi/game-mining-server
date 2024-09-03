package caches

import "strconv"

// GenRateLimitCacheKey generate rate limit cache key, limit:{biz}:{target}, e.g: limit:IP:127.0.0.1
func GenRateLimitCacheKey(mode string, target string) string {
	return "limit:" + mode + ":" + target
}

func GenCoinPriceCacheKey(fiatSymbol string, coinSymbol string) string {
	return "price:" + fiatSymbol + ":" + coinSymbol
}

// GenUserCacheKey generate user cache key, user:{uid}, e.g: user:102231405510
func GenUserCacheKey(uid int64) string {
	return "u:" + strconv.FormatInt(uid, 10)
}

// GenUserSessionCacheKey generate user cache key, session:{uid}, e.g: session:102231405510
func GenUserSessionCacheKey(uid int64) string {
	return "s:" + strconv.FormatInt(uid, 10)
}
