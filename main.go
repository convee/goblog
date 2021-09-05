package main

import (
	"blog/pkg/cache"
	"blog/pkg/config"
	"blog/pkg/elasticsearch"
	"blog/pkg/handler"
	"blog/pkg/mongo"
	"blog/pkg/mysql"
	"blog/pkg/xlog"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

func init() {
	var (
		confPath string
	)
	pflag.StringVarP(&confPath, "conf", "c", "./conf/dev.yml", "yaml conf path")
	pflag.Parse()
	viper.SetConfigFile(confPath)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("load config failed:", err)
		os.Exit(1)
	}
}
func main() {
	config.LoadConfig()
	xlog.Init()
	mysql.Init()
	cache.Init()
	mongo.Init()
	elasticsearch.Init()
	addr := viper.GetString("system.addr")
	log.Println("start serve: [", addr, "]")
	srv := &http.Server{
		Addr:    addr,
		Handler: handler.InitRouter(),
	}
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("server run:", err)
	}

}
