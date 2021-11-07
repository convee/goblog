package config

var (
	RedisConfig redisConfig
	MysqlConfig mysqlConfig
	EsConfig esConfig
)
func LoadConfig()  {
	loadLogConf()
	loadDbConf()
	loadRedisConf()
	loadMongoConf()
	loadEsConf()
}