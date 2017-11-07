package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/handlers"
	"github.com/tierpod/metatiles-cacher/pkg/logger"
	"github.com/tierpod/metatiles-cacher/pkg/queue"
)

var version string

func main() {
	// Command line flags
	var (
		flagVersion bool
		flagConfig  string
	)

	flag.BoolVar(&flagVersion, "version", false, "Show version and exit")
	flag.StringVar(&flagConfig, "config", "./config.json", "Path to config file")
	flag.Parse()

	if flagVersion {
		fmt.Printf("Version: %v\n", version)
		os.Exit(0)
	}

	cfg, err := config.Load(flagConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.New(os.Stdout, cfg.Log.Debug, cfg.Log.Datetime)

	cw, err := cache.NewFileCacheWriter(cfg.FileCache, logger)
	if err != nil {
		logger.Fatal(err)
	}

	uq := queue.NewUniq()

	http.Handle("/add/", handlers.LogConnection(
		addHandler{
			cache:  cw,
			queue:  uq,
			cfg:    cfg,
			logger: logger}, logger))
	http.Handle("/status", handlers.LogConnection(
		handlers.XToken(statusHandler{}, cfg.Writer.XToken, logger),
		logger))

	logger.Printf("Starting web server on %v", cfg.Writer.Bind)
	err = http.ListenAndServe(cfg.Writer.Bind, nil)
	if err != nil {
		logger.Fatal(err)
	}
}
