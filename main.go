package main

import (
	"log"
	"net/http"

	"github.com/convee/goblog/conf"
	"github.com/convee/goblog/internal/pkg/es"
	"github.com/convee/goblog/internal/pkg/mysql"
	"github.com/convee/goblog/internal/routers"
	logger "github.com/convee/goblog/pkg/log"
	"github.com/convee/goblog/pkg/redis"

	"github.com/spf13/pflag"
)

var (
	cfgFile = pflag.StringP("config", "c", "./conf/dev.yml", "config file path.")
	//version = pflag.BoolP("version", "v", false, "show version info.")
)

func main() {

	pflag.Parse()

	// init config
	cfg := conf.Init(*cfgFile)

	// init logger
	logger.Init(&cfg.Logger)

	// init redis
	redis.Init(&cfg.Redis)

	// init orm
	//model.Init(&cfg.ORM)

	// init mysql
	mysql.Init(&cfg.Mysql)

	// init elasticsearch
	if !cfg.Elasticsearch.Disable {
		es.Init(&cfg.Elasticsearch)
	}

	addr := cfg.App.Addr
	log.Println("start serve: [", addr, "]")
	srv := &http.Server{
		Addr:    addr,
		Handler: routers.InitRouter(),
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Println("server run:", err)
	}

}
