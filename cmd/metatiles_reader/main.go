// metatiles-reader provides http interface for serving tiles from metatiles cache.
// Contains slippy map based on LeafLet.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	// _ "net/http/pprof"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/config"
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

	cfg := config.NewConfig(flagConfig)

	logger := logger.NewLogger(os.Stdout, cfg.Reader.LogDebug, cfg.Reader.LogDatetime)

	filecache, err := cache.NewFileCacheReader(cfg.Reader.RootDir, logger)
	if err != nil {
		logger.Fatal(err)
	}

	http.Handle("/status", handlers.LogConnection(
		handlers.XToken(
			statusHandler{}, cfg.Reader.XToken, logger,
		),
		logger))
	http.Handle("/static/", handlers.LogConnection(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))), logger),
	)
	http.Handle("/maps/", handlers.LogConnection(
		mapsHandler{
			logger: logger,
			cache:  filecache,
			cfg:    cfg,
		}, logger))

	logger.Printf("Starting web server on: %v", cfg.Reader.Bind)
	err = http.ListenAndServe(cfg.Reader.Bind, nil)
	if err != nil {
		logger.Fatal(err)
	}
}
