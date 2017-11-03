package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/fetchservice"
	"github.com/tierpod/metatiles-cacher/pkg/handlers"
	"github.com/tierpod/metatiles-cacher/pkg/logger"
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

	cfg := config.Load(flagConfig)

	logger := logger.New(os.Stdout, cfg.Log.Debug, cfg.Log.Datetime)

	cw, err := cache.NewFileCacheWriter(cfg.FileCache, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// TODO: buffer?
	fetchservice := fetchservice.NewFetchService(1, cw, cfg, logger)

	http.Handle("/add/", handlers.LogConnection(
		addHandler{
			cache:  cw,
			fs:     fetchservice,
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
