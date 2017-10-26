package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/coords"
	"github.com/tierpod/metatiles-cacher/pkg/fetchservice"
	"github.com/tierpod/metatiles-cacher/pkg/httpclient"
)

type mapsHandler struct {
	logger     *log.Logger
	cache      cache.Reader
	sources    map[string]string
	writer     bool
	writerAddr string
}

func (h mapsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zxy, style, format, err := coords.NewZXYFromURL(r.URL.Path)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	source, found := h.sources[style]
	if !found {
		h.logger.Printf("[ERROR] Style not found in sources: %v", style)
		http.Error(w, "Style not found in sources", http.StatusNotFound)
		return
	}
	h.logger.Printf("[DEBUG] Convert URL(%v) to %v, style(%v), format(%v)", r.URL.Path, zxy, style, format)

	data, found, err := h.cache.Read(zxy, style)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !found {
		url := fmt.Sprintf(source, zxy.Path())
		h.logger.Printf("Fetch from upstream %v", url)
		data, err = httpclient.Get(url)
		if err != nil {
			h.logger.Printf("[ERROR] %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send request to metatiles-writer
		if h.writer {
			h.logger.Printf("Send request to writer: %v, style(%v)", zxy.ConvertToMeta(), style)
			job := fetchservice.NewJob(zxy.ConvertToMeta(), style, source)
			buf := new(bytes.Buffer)
			json.NewEncoder(buf).Encode(job)
			err := httpclient.PostJSON(h.writerAddr, buf)
			if err != nil {
				h.logger.Printf("[ERROR] %v\n", err)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
	return
}
