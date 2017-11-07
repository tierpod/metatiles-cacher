package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/coords"
	"github.com/tierpod/metatiles-cacher/pkg/httpclient"
	"github.com/tierpod/metatiles-cacher/pkg/utils"
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if zxy.Z < h.cfg.Zoom.Min || zxy.Z > h.cfg.Zoom.Max {
		h.logger.Printf("[ERROR] Wrong zoom level: Z(%v)", zxy.Z)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	source, found := h.cfg.Sources.Map[style]
	if !found {
		h.logger.Printf("[ERROR] Style not found in sources: %v", style)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.logger.Printf("[DEBUG] Got request %v style(%v), format(%v)", zxy, style, format)

	found, mtime := h.cache.Check(zxy, style)
	// found in cache
	if found {
		// check Etag
		etag := `"` + utils.DigestString(mtime.String()) + `"`

		w.Header().Set("Etag", etag)
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%v", h.cfg.Reader.MaxAge))

		if match := r.Header.Get("If-None-Match"); match != "" {
			if match == etag {
				h.logger.Printf("[DEBUG] File not modified: Etag(%v) == If-None-Match(%v)", etag, match)
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		data, readErr := h.cache.Read(zxy, style)
		if readErr != nil {
			h.logger.Printf("[ERROR] %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Write(data)
		return
	}

	// not found in cache
	url := strings.Replace(source, "{zxy}", zxy.Path(), 1)
	h.logger.Printf("Get from source %v", url)

	// get tile from remote source and reply to client
	if h.cfg.Reader.UseSources {
		data, err := httpclient.Get(url, h.cfg.HTTPClient.UserAgent)
		if err != nil {
			h.logger.Printf("[ERROR] %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Write(data)
	}

	// send writer request to metatiles_writer
	if h.cfg.Reader.UseWriter {
		go func() {
			url := h.cfg.Reader.WriterAddr + "/" + style + "/" + zxy.ConvertToMeta().Path()
			h.logger.Printf("Send request to writer: %v", url)
			_, err := httpclient.Get(url, h.cfg.HTTPClient.UserAgent)
			if err != nil {
				h.logger.Printf("[ERROR] %v\n", err)
				return
			}
		}()
	}

	return
}
