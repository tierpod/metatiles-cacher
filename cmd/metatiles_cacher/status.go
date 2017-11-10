package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/tierpod/metatiles-cacher/pkg/queue"
)

type statusHandler struct {
	queue *queue.Uniq
}

func (h statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Goroutines: %v\n", runtime.NumGoroutine())
	fmt.Fprintf(w, "Queue length: %v\n", h.queue.Len())
	fmt.Fprintf(w, "Queue items: %v\n", h.queue.Items())
	return
}
