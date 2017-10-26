package main

import (
	"fmt"
	"net/http"
	"runtime"
)

type statusHandler struct {
}

func (h statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Goroutines: %v\n", runtime.NumGoroutine())
}
