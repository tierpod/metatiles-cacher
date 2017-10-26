package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/fetchservice"
)

type addHandler struct {
	logger  *log.Logger
	cache   cache.Writer
	fs      *fetchservice.FetchService
	sources map[string]string
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

	_, found := h.sources[j.Style]
	if !found {
		h.logger.Printf("[ERROR] Style not found in sources: %v", j.Style)
		http.Error(w, "Style not found in sources", http.StatusNotFound)
		return
	}

	h.logger.Printf("Receive data: %v Style(%v)", j.Meta, j.Style)

	h.fs.Add(j)
}
