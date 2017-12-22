package main

import (
	"fmt"
	"net/http"
	"runtime"
)

// LockInfoer provides interface for get items locker.
type LockInfoer interface {
	Items() []string
}

type statusHandler struct {
	locker LockInfoer
}

func (h statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "goroutines: %v\n", runtime.NumGoroutine())
	fmt.Fprintf(w, "background locker items: %v\n", h.locker.Items())
	return
}
