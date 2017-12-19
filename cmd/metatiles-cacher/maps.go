package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/fetch"
	"github.com/tierpod/metatiles-cacher/pkg/latlong"
	"github.com/tierpod/metatiles-cacher/pkg/metatile"
	"github.com/tierpod/metatiles-cacher/pkg/tile"
	"github.com/tierpod/metatiles-cacher/pkg/util"
)

type mapsHandler struct {
	logger  *log.Logger
	cache   cache.ReadWriter
	cfg     *config.Config
	fetcher fetch.CacheWriter
}

func (h mapsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	minZoom := source.Zoom.Min
	maxZoom := source.Zoom.Max
	// if t.Zoom > maxZoom && source.HasRegion() { // only if tile max zoom always > region max zoom
	if source.HasRegion() {
		ll := latlong.New(t.Zoom, t.X, t.Y)
		if source.Region.Polygons.Contains(ll) {
			h.logger.Printf("[DEBUG] Point(%v) inside Region(%v)", t, source.Region.File)
			minZoom = source.Region.Zoom.Min
			maxZoom = source.Region.Zoom.Max
		}
	}

	if t.Zoom < minZoom || t.Zoom > maxZoom {
		h.logger.Printf("[ERROR] wrong zoom level for Source(%v): Zoom(%v)", source.Name, t.Zoom)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	mimetype, err := util.Mimetype(t.Ext)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.logger.Printf("[DEBUG] try get tile from cache")
	found, mtime := h.cache.Check(t)
	if found {
		etag := `"` + util.DigestString(mtime.String()) + `"`
		h.replyFromCache(w, t, mimetype, etag, r.Header.Get("If-None-Match"))
		return
	}

	// fetch tiles for metatile and write to cache?
	mt := metatile.NewFromTile(t)
	err = h.fetcher.MetatileWriteToCache(mt, source.URL, h.cache)
	if err != nil {
		h.logger.Printf("[ERROR]: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Printf("[DEBUG] try get tile from cache after writing")
	// try again
	found, mtime = h.cache.Check(t)
	if found {
		etag := `"` + util.DigestString(mtime.String()) + `"`
		h.replyFromCache(w, t, mimetype, etag, r.Header.Get("If-None-Match"))
		return
	}

	h.logger.Printf("[ERROR] unable to get tile")
	w.WriteHeader(http.StatusNotFound)
	return
}

func (h mapsHandler) replyFromCache(w http.ResponseWriter, t tile.Tile, mimetype, etag, ifNoneMatch string) {
	w.Header().Set("Etag", etag)
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%v", h.cfg.Service.MaxAge))

	if ifNoneMatch == etag {
		h.logger.Printf("[DEBUG] replyFromCache: file not modified: Etag(%v) == If-None-Match(%v)", etag, ifNoneMatch)
		w.WriteHeader(http.StatusNotModified)
		return
	}

	data, err := h.cache.Read(t)
	if err != nil {
		h.logger.Printf("[ERROR] replyFromCache: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", mimetype)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
	return
}
