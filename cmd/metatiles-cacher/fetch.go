package main

import (
	"log"
	"net/http"

	"github.com/tierpod/go-osm/metatile"
	"github.com/tierpod/go-osm/tile"

	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/fetch"
)

type fetchHandler struct {
	logger *log.Logger
	cfg    *config.Config
	fs     *fetch.Service
}

func (h fetchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, err := tile.NewFromURL(r.URL.Path)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	mt := metatile.NewFromTile(t)

	source, err := h.cfg.Source(mt.Style)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		return
	}

	h.logger.Printf("[DEBUG] got request: %v", mt)
	err = h.fs.AddWait(mt, source.URL)
	if err != nil {
		if err == fetch.ErrJobInQueue {
			w.WriteHeader(http.StatusAccepted)
			h.logger.Printf("[WARN] %v already in queue, skip", t)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Printf("[ERROR] %v", err)
		return
	}

	return
}
