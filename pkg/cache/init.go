package cache

import (
	config2 "blog/pkg/config"
	"github.com/garyburd/redigo/redis"
	"time"
)

var (
	Redis *redis.Pool
)

func Init() {
	Redis = &redis.Pool{
		MaxIdle:     config2.RedisConfig.MaxIdle,                                  //最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
		MaxActive:   config2.RedisConfig.MaxActive,                                //最大的激活连接数，表示同时最多有N个连接
		IdleTimeout: time.Duration(config2.RedisConfig.IdleTimeout) * time.Second, //最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		Dial: func() (redis.Conn, error) { //建立连接
			return redis.Dial("tcp", config2.RedisConfig.Address)
		},
	}
}
