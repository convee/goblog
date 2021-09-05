package config

import "github.com/spf13/viper"

type redisConfig struct {
	Address     string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int64
}

func loadRedisConf()  {
	RedisConfig.Address = viper.GetString("redis.addr")
	RedisConfig.MaxIdle = viper.GetInt("redis.max_idle")
	RedisConfig.MaxActive = viper.GetInt("redis.max_active")
	RedisConfig.IdleTimeout = viper.GetInt64("redis.idle_timeout")
}
