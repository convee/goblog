package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/convee/goblog/conf"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/routers"
	"github.com/convee/goblog/internal/security"
	"github.com/convee/goblog/pkg/logger"
	"github.com/spf13/pflag"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
	cfgFile   = pflag.StringP("config", "c", conf.DefaultPath(), "config file path")
	showVer   = pflag.BoolP("version", "v", false, "print version information")
)

func main() {
	pflag.Parse()
	if *showVer {
		fmt.Printf("goblog %s (commit %s, built %s)\n", version, commit, buildTime)
		return
	}
	cfg := conf.Init(*cfgFile)
	logger.Init(&cfg.Logger)
	if err := security.ConfigureTrustedProxies(cfg.App.TrustedProxies); err != nil {
		log.Fatalf("configure trusted proxies: %v", err)
	}
	secret := os.Getenv("GOBLOG_SESSION_SECRET")
	if secret == "" {
		if !strings.EqualFold(cfg.App.Mode, "debug") {
			log.Fatal("GOBLOG_SESSION_SECRET is required outside debug mode")
		}
		log.Print("warning: using an ephemeral session secret in debug mode; sessions will expire after restart")
	} else if err := security.ConfigureSessionSecret(secret); err != nil {
		log.Fatalf("configure session secret: %v", err)
	}

	mysqlDB := daos.Init(&cfg.Mysql)

	srv := &http.Server{
		Addr:              cfg.App.Addr,
		Handler:           routers.InitRouter(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", cfg.App.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	select {
	case <-signalCtx.Done():
		log.Println("shutdown signal received")
	case err := <-serverErrors:
		log.Printf("server stopped unexpectedly: %v", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http server shutdown: %v", err)
	}
	if err := mysqlDB.Close(); err != nil {
		log.Printf("mysql close: %v", err)
	}
	if signalCtx.Err() == nil {
		os.Exit(1)
	}
}
