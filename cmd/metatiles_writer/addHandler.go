package main

import (
	"encoding/json"
	"log"
	"net/http"

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
	var j fetchservice.Job

	if r.Body == nil {
		h.logger.Printf("[ERROR] Empty Body")
		http.Error(w, "Empty Body", http.StatusNotFound)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		h.logger.Printf("[ERROR]: Json decode: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if j.Meta.Z < h.cfg.Writer.MinZoom || j.Meta.Z > h.cfg.Writer.MaxZoom {
		h.logger.Printf("[ERROR] Wrong zoom level: Z(%v)", j.Meta.Z)
		http.Error(w, "Wrong zoom level", http.StatusNotFound)
		return
	}

	source, found := h.cfg.SourcesMap[j.Style]
	if !found {
		h.logger.Printf("[ERROR] Style not found in sources: %v", j.Style)
		http.Error(w, "Style not found in sources", http.StatusNotFound)
		return
	}
	j.Source = source

	h.logger.Printf("Receive data: %+v", j)

	h.fs.Add(j)
}
