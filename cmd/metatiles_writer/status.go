package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

type statusHandler struct {
	logger *log.Logger
}

func (h statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Goroutines: %v\n", runtime.NumGoroutine())
}
