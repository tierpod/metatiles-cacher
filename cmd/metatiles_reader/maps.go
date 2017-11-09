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
	"github.com/tierpod/metatiles-cacher/pkg/util"
)

type mapsHandler struct {
	logger *log.Logger
	cache  cache.Reader
	cfg    *config.Service
}

func (h mapsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, style, format, err := coords.NewTileFromURL(r.URL.Path)
	if err != nil {
		h.logger.Printf("[ERROR] Wrong request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if t.Zoom < h.cfg.Zoom.Min || t.Zoom > h.cfg.Zoom.Max {
		h.logger.Printf("[ERROR] Wrong zoom level: Zoom(%v)", t.Zoom)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	source, found := h.cfg.Sources.Map[style]
	if !found {
		h.logger.Printf("[ERROR] Style not found in sources: %v", style)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.logger.Printf("[DEBUG] Got request %v style(%v), format(%v)", t, style, format)

	found, mtime := h.cache.Check(t, style)
	// found in cache
	if found {
		etag := `"` + util.DigestString(mtime.String()) + `"`
		h.replyFromCache(w, t, style, etag, r.Header.Get("If-None-Match"))
		return
	}

	// not found in cache
	if h.cfg.Reader.UseSources {
		h.replyFromSource(w, t, source)
	}

	// send request to writer
	if h.cfg.Reader.UseWriter {
		go h.sendToWriter(w, t, style)
	}

	return
}

func (h mapsHandler) replyFromCache(w http.ResponseWriter, t coords.Tile, style string, etag, ifNoneMatch string) {
	w.Header().Set("Etag", etag)
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%v", h.cfg.Reader.MaxAge))

	if ifNoneMatch == etag {
		h.logger.Printf("[DEBUG] File not modified: Etag(%v) == If-None-Match(%v)", etag, ifNoneMatch)
		w.WriteHeader(http.StatusNotModified)
		return
	}

	data, err := h.cache.Read(t, style)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
	return
}

func (h mapsHandler) replyFromSource(w http.ResponseWriter, t coords.Tile, source string) {
	url := strings.Replace(source, "{tile}", t.Path(), 1)
	h.logger.Printf("Get from source %v", url)

	data, err := httpclient.Get(url, h.cfg.HTTPClient.UserAgent)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
	return
}

func (h mapsHandler) sendToWriter(w http.ResponseWriter, t coords.Tile, style string) {
	url := h.cfg.Reader.WriterAddr
	url = strings.Replace(url, "{style}", style, 1)
	url = strings.Replace(url, "{metatile}", t.ToMetatile().Path(), 1)
	h.logger.Printf("Send request to writer: %v", url)

	_, err := httpclient.Get(url, h.cfg.HTTPClient.UserAgent)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		return
	}

	return
}
