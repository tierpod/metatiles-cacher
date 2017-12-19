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
	fetcher fetch.CacheWriter
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
	err = h.fetcher.MetatileWriteToCache(mt, source.URL, h.cache)
	if err != nil {
		if err == fetch.ErrQueueHasKey {
			w.WriteHeader(http.StatusCreated)
			return
		}

		h.logger.Printf("[ERROR]: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}
