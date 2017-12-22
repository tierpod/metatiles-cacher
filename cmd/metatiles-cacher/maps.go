package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/latlong"
	"github.com/tierpod/metatiles-cacher/pkg/metatile"
	"github.com/tierpod/metatiles-cacher/pkg/tile"
	"github.com/tierpod/metatiles-cacher/pkg/util"
)

// MetatileCacheReadWriter is the interface for using metatile cache.
type MetatileCacheReadWriter interface {
	Read(mc cache.MetatileCacher, t tile.Tile, w io.Writer) error
	Write(mc cache.MetatileCacher, r io.Reader) error
	Check(mc cache.MetatileCacher) (mtime time.Time, found bool)
}

// LockWaiter is the interface for using locks and wait to unlocks.
type LockWaiter interface {
	Add(key string)
	Del(key string)
	Wait(key string)
	HasKey(key string) bool
}

type mapsHandler struct {
	logger *log.Logger
	cacher MetatileCacheReadWriter
	cfg    *config.Config
	locker LockWaiter
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

	mt := metatile.NewFromTile(t)

	// if key found in locker struct, wait while Del will be executed
	if h.locker.HasKey(mt.Filepath("")) {
		h.logger.Printf("[DEBUG] wait for another handler complete fetching")
		h.locker.Wait(mt.Filepath(""))
	}

	h.logger.Printf("[DEBUG] try get tile from cache")
	mtime, found := h.cacher.Check(mt)
	if found {
		etag := `"` + util.DigestString(mtime.String()) + `"`
		h.replyFromCache(w, t, mimetype, etag, r.Header.Get("If-None-Match"))
		return
	}

	// before starting long-running fetching, add metatile info to locker to prevent another fetching
	// for this metatile. After long-running fetching, Del metatile info from locker (and notify
	// another waiters).
	h.locker.Add(mt.Filepath(""))
	defer h.locker.Del(mt.Filepath(""))

	// fetch metatile to buffer
	var buf bytes.Buffer
	err = mt.FetchDataEncodeTo(&buf, source.URL, h.cfg.Fetch.UserAgent)
	if err != nil {
		h.logger.Printf("[ERROR]: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write metatile to file
	err = h.cacher.Write(mt, &buf)
	if err != nil {
		h.logger.Printf("[ERROR]: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Printf("[DEBUG] try get tile from cache after writing")
	// try again
	mtime, found = h.cacher.Check(mt)
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

	mt := metatile.NewFromTile(t)
	var buf bytes.Buffer
	err := h.cacher.Read(mt, t, &buf)
	if err != nil {
		h.logger.Printf("[ERROR] replyFromCache: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", mimetype)
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	io.Copy(w, &buf)
	return
}
