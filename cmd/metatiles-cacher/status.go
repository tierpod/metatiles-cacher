package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/fetch"
)

type statusHandler struct {
	logger *log.Logger
	fs     *fetch.Service
	cfg    *config.Config
}

func (h statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "goroutines: %v\n", runtime.NumGoroutine())

	fmt.Fprintf(w, "last update time:\n")
	for name, source := range h.cfg.Sources {
		fmt.Fprintf(w, " %v: %v\n", name, source.LastUpdateTime)
	}

	jobs := h.fs.Jobs()
	fmt.Fprintf(w, "jobs in progress (%v):\n", len(jobs))
	for _, j := range jobs {
		fmt.Fprintf(w, " %v\n", j)
	}

	return
}
