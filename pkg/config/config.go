package config

var (
	RedisConfig redisConfig
	MysqlConfig mysqlConfig
)
func LoadConfig()  {
	loadDbConf()
	loadRedisConf()
	loadMongoConf()
	loadLogConf()
}