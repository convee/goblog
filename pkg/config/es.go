package config

import (
	"github.com/spf13/viper"
	"strings"
)

type esConfig struct {
	Urls     []string
}

func loadEsConf()  {
	EsConfig.Urls = strings.Split(viper.GetString("es.urls"), ",")
}