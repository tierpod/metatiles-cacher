// metatiles-reader provides http interface for serving tiles from metatiles cache.
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

	cfg, err := config.Load(flagConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.New(os.Stdout, cfg.Log.Debug, cfg.Log.Datetime)

	cr, err := cache.NewFileCacheReader(cfg.FileCache, logger)
	if err != nil {
		logger.Fatal(err)
	}

	http.Handle("/status", handler.LogConnection(
		handler.XToken(
			statusHandler{}, cfg.Reader.XToken, logger,
		),
		logger))
	http.Handle("/static/", handler.LogConnection(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))), logger),
	)
	http.Handle("/maps/", handler.LogConnection(
		mapsHandler{
			logger: logger,
			cache:  cr,
			cfg:    cfg,
		}, logger))

	logger.Printf("Starting web server on: %v", cfg.Reader.Bind)
	err = http.ListenAndServe(cfg.Reader.Bind, nil)
	if err != nil {
		logger.Fatal(err)
	}
}
