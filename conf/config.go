package conf

import (
	"github.com/convee/goblog/pkg/logger"
	"github.com/convee/goblog/pkg/redis"
	"github.com/convee/goblog/pkg/storage/elasticsearch"
	"github.com/convee/goblog/pkg/storage/mysql"
	"github.com/convee/goblog/pkg/storage/orm"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config global config
type Config struct {
	// common
	App AppConfig
	// component config
	Logger        logger.Config
	ORM           orm.Config
	Mysql         mysql.Config
	Redis         redis.Config
	Elasticsearch elasticsearch.Config
}

// AppConfig app config
type AppConfig struct {
	Name            string
	Version         bool
	Mode            string
	Addr            string
	Host            string
	Cdn             string
	DisableDingDing bool
}

var (
	// Conf app global config
	Conf = &Config{}
)

func Init(configPath string) *Config {
	viper.SetConfigType("yml")
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(&Conf); err != nil {
			panic(err)
		}
	})
	return Conf
}
