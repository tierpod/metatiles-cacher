package main

import (
	"log"
	"net/http"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/fetch"
	"github.com/tierpod/metatiles-cacher/pkg/metatile"
	"github.com/tierpod/metatiles-cacher/pkg/tile"
	"github.com/tierpod/metatiles-cacher/pkg/util"
)

type fetchHandler struct {
	logger  *log.Logger
	cache   cache.ReadWriter
	cfg     *config.Config
	fetcher *fetch.Fetch
}

func (h fetchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, err := tile.NewFromURL(r.URL.Path)
	if err != nil {
		h.logger.Printf("[ERROR] wrong request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Printf("[DEBUG] got request %v", t)

	source, err := h.cfg.Source(t.Map)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if t.Zoom < source.Zoom.Min || t.Zoom > source.Zoom.Max {
		h.logger.Printf("[ERROR] wrong zoom level for Source(%v): Zoom(%v)", source.Name, t.Zoom)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	_, err = util.Mimetype(t.Ext)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// fetch tiles for metatile and write to cache
	mt := metatile.NewFromTile(t)
	data, err := h.fetcher.Metatile(mt, source.URL)
	if err != nil && err.Error() == "item already present in queue" {
		w.WriteHeader(http.StatusCreated)
		return
	} else if err != nil {
		h.logger.Printf("[ERROR]: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.cache.Write(mt, data)
	if err != nil {
		h.logger.Printf("[ERROR]: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
