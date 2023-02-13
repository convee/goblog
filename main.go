package main

import (
	"context"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/es"
	"github.com/convee/goblog/pkg/logger"
	"log"
	"net/http"
	"time"

	"github.com/convee/goblog/conf"
	"github.com/convee/goblog/internal/routers"
	"github.com/convee/goblog/pkg/redis"
	"github.com/convee/goblog/pkg/shutdown"

	"github.com/spf13/pflag"
)

var (
	cfgFile = pflag.StringP("config", "c", "./conf/local.yml", "config file path.")
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
	daos.Init(&cfg.Mysql)

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
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("server run:", err)
		}
	}()

	shutdown.NewHook().Close(
		// 关闭 http server
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Println("http server closed err", err)
			} else {
				log.Println("http server closed")
			}
		},
	)

}
