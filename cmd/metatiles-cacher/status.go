package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/tierpod/metatiles-cacher/pkg/fetch"
)

type statusHandler struct {
	logger *log.Logger
	fs     *fetch.Service
}

func (h statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "goroutines: %v\n", runtime.NumGoroutine())

	jobs := h.fs.Jobs()
	fmt.Fprintf(w, "Jobs in progress (%v):\n", len(jobs))
	fmt.Fprintf(w, "%v\n", jobs)
	return
}
