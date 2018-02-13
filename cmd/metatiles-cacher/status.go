package main

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/tierpod/go-osm/metatile"

	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/fetch"
)

type statusHandler struct {
	logger *log.Logger
	fs     *fetch.Service
	cfg    *config.Config
}

type statusResult struct {
	Goroutines     int                  `json:"goroutines"`
	LastUpdateTime map[string]time.Time `json:"last_update_time"`
	Jobs           []metatile.Metatile  `json:"jobs"`
}

func (h statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lu := make(map[string]time.Time, len(h.cfg.Sources))
	for name, source := range h.cfg.Sources {
		lu[name] = source.LastUpdateTime
	}

	status := statusResult{
		Goroutines:     runtime.NumGoroutine(),
		LastUpdateTime: lu,
		Jobs:           h.fs.Jobs(),
	}

	result, err := json.Marshal(status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Printf("[ERROR] status result: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	return
}
