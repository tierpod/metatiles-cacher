package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/coords"
	"github.com/tierpod/metatiles-cacher/pkg/httpclient"
)

type mapsHandler struct {
	logger *log.Logger
	cache  cache.Reader
	cfg    *config.Service
}

func (h mapsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zxy, style, format, err := coords.NewZXYFromURL(r.URL.Path)
	if err != nil {
		h.logger.Printf("[ERROR] Wrong request: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if zxy.Z < h.cfg.Reader.MinZoom || zxy.Z > h.cfg.Reader.MaxZoom {
		h.logger.Printf("[ERROR] Wrong zoom level: Z(%v)", zxy.Z)
		http.Error(w, "Wrong zoom level", http.StatusNotFound)
		return
	}

	source, found := h.cfg.SourcesMap[style]
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

	// if found in cache
	if found {
		mapsReply(w, data)
		return
	}

	url := strings.Replace(source, "{zxy}", zxy.Path(), 1)
	h.logger.Printf("Get from source %v", url)

	data, err = httpclient.Get(url)
	if err != nil {
		h.logger.Printf("[ERROR] %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mapsReply(w, data)

	// Send request to metatiles_writer
	if h.cfg.Reader.WriterAddr != "" {
		go func() {
			h.logger.Printf("Send request to writer: %v, style(%v)", zxy.ConvertToMeta(), style)
			url := h.cfg.Reader.WriterAddr + "/" + style + "/" + zxy.ConvertToMeta().Path()
			_, err := httpclient.Get(url)
			if err != nil {
				h.logger.Printf("[ERROR] %v\n", err)
				return
			}
		}()
	}

	return
}

func mapsReply(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
}
