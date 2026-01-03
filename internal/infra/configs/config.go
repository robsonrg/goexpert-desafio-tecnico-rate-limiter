package configs

import (
	"path/filepath"

	"github.com/spf13/viper"
)

var cfg *conf

type conf struct {
	RedisHost                       string `mapstructure:"REDIS_HOST"`
	RedisPassword                   string `mapstructure:"REDIS_PASSWORD"`
	RateLimitRequestPerSecondsIP    int64  `mapstructure:"RATE_LIMIT_REQUESTS_PER_SECOND_IP"`
	RateLimitRequestPerSecondsToken int64  `mapstructure:"RATE_LIMIT_REQUESTS_PER_SECOND_TOKEN"`
}

func LoadConfig(path string) (*conf, error) {
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(filepath.Join(path, ".env"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, nil
}
