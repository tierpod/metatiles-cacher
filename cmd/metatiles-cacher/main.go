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

	"github.com/gorilla/handlers"

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
	tokens := newTokenStore(cfg.HTTP.XToken, logger)

	r := http.NewServeMux()

	// admin handlers
	r.Handle("/status", tokens.Middleware(
		statusHandler{
			logger: logger,
			fs:     fs,
			cfg:    cfg,
		},
	))

	// static handler
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// service handlers
	r.Handle("/maps/", mapsHandler{
		cfg:    cfg,
		logger: logger,
		cache:  metaCache,
		fs:     fs,
	})
	r.Handle("/fetch/", fetchHandler{
		logger: logger,
		cfg:    cfg,
		fs:     fs,
	})

	logger.Printf("Starting web server on: %v", cfg.HTTP.Bind)
	err = http.ListenAndServe(cfg.HTTP.Bind, handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		logger.Fatal(err)
	}
}
