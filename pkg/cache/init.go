package cache

import (
	config2 "blog/pkg/config"
	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
)

func Init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:	  config2.RedisConfig.Address,
		DB:		  0,
	})
}
