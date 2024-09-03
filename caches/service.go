package caches

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"game-mining-server/configs"
	"game-mining-server/utils"
	"github.com/go-redis/redis_rate/v10"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type Service struct {
	RdsInstance *redis.Client       // Redis client
	RateLimiter *redis_rate.Limiter // Api rate limiter
	RedSyncLock *redsync.Redsync    // Multi lock
}

func CreateCacheService(cfg *configs.CacheConfig) *Service {
	rdb := redis.NewClient(&redis.Options{
		Addr:      fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:  cfg.Pass,
		DB:        0,
		PoolSize:  10,
		TLSConfig: utils.Any(cfg.Pass == "", &tls.Config{MinVersion: tls.VersionTLS12}, nil),
	})
	dbSize := rdb.DBSize(context.Background())
	log.Printf("Create cache service, db %s", dbSize.String())
	return &Service{RdsInstance: rdb, RateLimiter: redis_rate.NewLimiter(rdb), RedSyncLock: redsync.New(goredis.NewPool(rdb))}
}

// SetString store a key-value pair to Cache service
func (s *Service) SetString(key string, item string, expiresSec int) error {
	return s.RdsInstance.Set(context.Background(), key, item, time.Duration(expiresSec)*time.Second).Err()
}

// GetString try to get a string value by key from Cache service
func (s *Service) GetString(key string) (string, error) {
	return s.RdsInstance.Get(context.Background(), key).Result()
}

// Delete try to delete al value in Cache service
func (s *Service) Delete(key string) {
	s.RdsInstance.Del(context.Background(), key)
}

// HSetStruct Set a struct data into redis hash table
func (s *Service) HSetStruct(key string, data interface{}, expiresSec int) error {
	if expiresSec <= 0 {
		return s.RdsInstance.HSet(context.Background(), key, data).Err()
	} else {
		_, err := s.RdsInstance.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
			if e0 := pipe.HSet(context.Background(), key, data).Err(); e0 != nil {
				return e0
			}
			if e1 := pipe.Expire(context.Background(), key, time.Duration(expiresSec)*time.Second).Err(); e1 != nil {
				return e1
			}
			return nil
		})
		return err
	}
}

// HGetStruct Get a struct data from redis hash table
func (s *Service) HGetStruct(key string, v interface{}) error {
	result := s.RdsInstance.HGetAll(context.Background(), key)
	fields, err := result.Result()
	if err != nil {
		return err
	}
	if len(fields) == 0 {
		return errors.New("HGet not found")
	}
	return result.Scan(v)
}

// HMapSet Set a map to cache
func (s *Service) HMapSet(key string, m map[string]interface{}, expiresSec int) error {
	if expiresSec <= 0 {
		return s.RdsInstance.HMSet(context.Background(), key, m).Err()
	} else {
		_, err := s.RdsInstance.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
			if e0 := pipe.HMSet(context.Background(), key, m).Err(); e0 != nil {
				return e0
			}
			if e1 := pipe.Expire(context.Background(), key, time.Duration(expiresSec)*time.Second).Err(); e1 != nil {
				return e1
			}
			return nil
		})
		return err
	}
}

// HMapGet Get a map data
func (s *Service) HMapGet(key string) (map[string]string, error) {
	if fields, err := s.RdsInstance.HGetAll(context.Background(), key).Result(); err != nil {
		return nil, err
	} else {
		return fields, nil
	}
}

// SSetAdd Redis SADD to add a value in key set
func (s *Service) SSetAdd(key string, value interface{}) error {
	if _, e := s.RdsInstance.SAdd(context.Background(), key, value).Result(); e != nil {
		return e
	} else {
		return nil
	}
}

// SSetDel Redis SREM to delete a value in key set
func (s *Service) SSetDel(key string, value interface{}) error {
	if _, e := s.RdsInstance.SRem(context.Background(), key, value).Result(); e != nil {
		return e
	} else {
		return nil
	}
}

// SSetCount Redis SCARD to query count in a set
func (s *Service) SSetCount(key string) (int64, error) {
	return s.RdsInstance.SCard(context.Background(), key).Result()
}
