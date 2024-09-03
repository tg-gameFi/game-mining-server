package app

import (
	"fmt"
	"game-mining-server/caches"
	"game-mining-server/configs"
	"game-mining-server/dbs"
	"github.com/mymmrac/telego"
)

type App struct {
	Config *configs.Config
	Bot    *telego.Bot
	DB     *dbs.Service
	Cache  *caches.Service
}

var instance App

func CreateApp(cfgPath string, botToken string) error {
	cfg := configs.LoadConfig[configs.Config](cfgPath)
	bot, e0 := createBot(botToken)
	if e0 != nil {
		return e0
	}

	db := dbs.CreateDBService(cfg.Database, cfg.Basic.Env)
	instance = App{
		Config: cfg,
		Bot:    bot,
		DB:     db,
		Cache:  caches.CreateCacheService(cfg.Cache),
	}
	return nil
}

func createBot(botToken string) (*telego.Bot, error) {
	if botToken == "" {
		return nil, nil
	}
	bot, e0 := telego.NewBot(botToken)
	if e0 != nil {
		return nil, fmt.Errorf("Create bot failed: " + e0.Error())
	}

	return bot, nil
}

func Config() *configs.Config {
	return instance.Config
}

func Cache() *caches.Service {
	return instance.Cache
}

func DB() *dbs.Service {
	return instance.DB
}

func Bot() *telego.Bot {
	return instance.Bot
}
