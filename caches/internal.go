package caches

import (
	"game-mining-server/dbs"
)

// UserFindByIdCached Find User in cache and if there is no cache found, find in DB and update cache
func UserFindByIdCached(cacheService *Service, dbService *dbs.Service, uid int64, expiresSec int) (*dbs.User, error) {
	key := GenUserCacheKey(uid)
	var userCache dbs.User
	if e1 := cacheService.HGetStruct(key, &userCache); e1 != nil || &userCache == nil { // not found user in cache
		if userDb, e2 := dbService.UserFindById(uid); e2 != nil || userDb == nil {
			return nil, e2
		} else {
			_ = cacheService.HSetStruct(key, userDb, expiresSec)
			return userDb, nil
		}
	} else {
		return &userCache, nil
	}
}
