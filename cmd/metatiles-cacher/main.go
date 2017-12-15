// metatiles-cacher is the small web service for serving tiles from metatiles cache. If tile not
// found in cache, get it from remote source and write to metatiles cache.
//
// Contains slippy map based on LeafLet.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	// _ "net/http/pprof"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/handler"
	"github.com/tierpod/metatiles-cacher/pkg/httpclient"
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
	flag.StringVar(&flagConfig, "config", "./config.yaml", "Path to config file")
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

	fc, err := cache.NewFileCache(cfg.FileCache, logger)
	if err != nil {
		logger.Fatal(err)
	}

	fetcher := httpclient.NewFetch(cfg.HTTPClient, logger)

	uq := queue.NewUniq()

	http.Handle("/status", handler.LogConnection(
		handler.XToken(
			statusHandler{queue: uq}, cfg.Service.XToken, logger,
		),
		logger))
	http.Handle("/static/", handler.LogConnection(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))), logger),
	)
	http.Handle("/maps/", handler.LogConnection(
		mapsHandler{
			logger:  logger,
			cache:   fc,
			cfg:     cfg,
			fetcher: fetcher,
		}, logger))
	http.Handle("/fetch/", handler.LogConnection(
		fetchHandler{
			logger:  logger,
			cache:   fc,
			cfg:     cfg,
			fetcher: fetcher,
		}, logger))

	logger.Printf("Starting web server on: %v", cfg.Service.Bind)
	err = http.ListenAndServe(cfg.Service.Bind, nil)
	if err != nil {
		logger.Fatal(err)
	}
}
