package config

import "github.com/spf13/viper"

type redisConfig struct {
	Address     string
}

func loadRedisConf()  {
	RedisConfig.Address = viper.GetString("redis.addr")
}
