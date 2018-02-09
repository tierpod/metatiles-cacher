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
	"github.com/tierpod/metatiles-cacher/pkg/fetch"
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
	flag.StringVar(&flagConfig, "config", "./config.yaml", "Path to config file")
	flag.Parse()

	if flagVersion {
		fmt.Printf("Version: %v\n", version)
		os.Exit(0)
	}

	// configure service
	cfg, err := config.Load(flagConfig)
	if err != nil {
		log.Fatal(err)
	}

	// init logger
	logger := logger.New(os.Stdout, cfg.Log.Debug, cfg.Log.Datetime)

	// init metatiles cache
	metaCache, err := cache.NewMetatileCache(cfg.Cache, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// init and start background fetch service
	fs := fetch.NewService(cfg, metaCache, logger)
	fs.Start()

	// init token store for admin entrypoints
	// tokens = newTokenStore(cfg.HTTP.XToken, logger)

	// http.Handle("/status", handler.LogConnection(
	// 	handler.XToken(
	// 		statusHandler{locker: locker}, cfg.Service.XToken, logger,
	// 	),
	// 	logger))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/maps/", mapsHandler{
		cfg:    cfg,
		logger: logger,
		cache:  metaCache,
		fs:     fs,
	})
	http.Handle("/fetch/", fetchHandler{
		logger: logger,
		cfg:    cfg,
		fs:     fs,
	})

	logger.Printf("Starting web server on: %v", cfg.HTTP.Bind)
	err = http.ListenAndServe(cfg.HTTP.Bind, nil)
	if err != nil {
		logger.Fatal(err)
	}
}
