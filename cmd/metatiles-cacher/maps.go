package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/tierpod/go-osm/point"
	"github.com/tierpod/go-osm/tile"
	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/util"
)

// CacheReader provides interface for reading tiles from cache.
type CacheReader interface {
	Read(t tile.Tile) ([]byte, error)
	Check(t tile.Tile) (mtime time.Time, found bool)
}

type mapsHandler struct {
	logger *log.Logger
	cache  CacheReader
	cfg    *config.Config
}

func (h mapsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// parse incoming request
	t, err := tile.NewFromURL(r.URL.Path)
	if err != nil {
		h.logger.Printf("[ERROR] wrong request string: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Printf("[DEBUG] got request %v", t)

	// validate
	source, err := h.cfg.Source(t.Style)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	minZoom := config.MinZoom
	maxZoom := source.MaxZoom
	if t.Zoom > maxZoom && source.HasRegion() {
		p := point.ZXY{Z: t.Zoom, X: t.X, Y: t.Y}
		ll := p.ToLatLong()
		if source.Region.Polygons.Contains(ll) {
			h.logger.Printf("[DEBUG] Point(%v) inside Region(%v)", t, source.Region.File)
			maxZoom = source.Region.MaxZoom
		}
	}

	if t.Zoom < minZoom || t.Zoom > maxZoom {
		h.logger.Printf("[ERROR] forbidden zoom level (%v) for Source(%v)", t.Zoom, t.Style)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	mimetype, err := util.Mimetype(t.Ext)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.logger.Printf("[DEBUG] %v: try to get from cache", t)
	mtime, found := h.cache.Check(t)
	if found {
		etag := `"` + util.DigestString(mtime.String()) + `"`
		h.replyFromCache(w, t, mimetype, etag, r.Header.Get("If-None-Match"))
		return
	}

	h.logger.Printf("[DEBUG] %v: not found in cache, get from remote source", t)
	//h.replyFromSource()
	//h.sendToFetchService()
	return
}

func (h mapsHandler) replyFromCache(w http.ResponseWriter, t tile.Tile, mimetype, etag, ifNoneMatch string) {
	w.Header().Set("Etag", etag)

	if ifNoneMatch == etag {
		h.logger.Printf("[DEBUG] cache: file not modified: Etag(%v) == If-None-Match(%v)", etag, ifNoneMatch)
		w.WriteHeader(http.StatusNotModified)
		return
	}

	data, err := h.cache.Read(t)
	if err != nil {
		h.logger.Printf("[ERROR] cache: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for k, v := range h.cfg.HTTP.Headers {
		w.Header().Set(k, v)
	}
	w.Header().Set("Content-Type", mimetype)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
}
