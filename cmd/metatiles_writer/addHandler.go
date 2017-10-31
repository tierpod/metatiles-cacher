package main

import (
	"log"
	"net/http"

	"github.com/tierpod/metatiles-cacher/pkg/coords"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/fetchservice"
)

type addHandler struct {
	logger *log.Logger
	cache  cache.Writer
	fs     *fetchservice.FetchService
	cfg    *config.Service
}

func (h addHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	meta, style, err := coords.NewMetaFromURL(r.URL.Path)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if meta.Z < h.cfg.Writer.MinZoom || meta.Z > h.cfg.Writer.MaxZoom {
		h.logger.Printf("[ERROR] Wrong zoom level: Z(%v)", meta.Z)
		http.Error(w, "Wrong zoom level", http.StatusNotFound)
		return
	}

	source, found := h.cfg.SourcesMap[style]
	if !found {
		h.logger.Printf("[ERROR] Style not found in sources: %v", style)
		http.Error(w, "Style not found in sources", http.StatusNotFound)
		return
	}

	j := fetchservice.NewJob(meta, style, source)

	h.logger.Printf("Receive data: %+v", j)

	h.fs.Add(j)
}
