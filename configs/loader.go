package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

func LoadConfig[T any](path string) *T {
	viper.SetConfigType("json")
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("[Config] Load configs from %s failed: %s", path, err))
	}

	var cfg *T
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("[Config] Load configs from %s unmarshal failed: %s", path, err))
	}
	return cfg
}
