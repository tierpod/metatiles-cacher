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
	"github.com/tierpod/metatiles-cacher/pkg/queue"
	"github.com/tierpod/metatiles-cacher/pkg/util"
)

type mapsHandler struct {
	logger *log.Logger
	cache  cache.ReadWriter
	cfg    *config.Config

	queue *queue.Uniq
}

func (h mapsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, style, err := coords.NewTileFromURL(r.URL.Path)
	if err != nil {
		h.logger.Printf("[ERROR] Wrong request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Printf("[DEBUG] Got request %v style(%v)", t, style)

	source, err := h.cfg.Source(style)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	minZoom := source.Zoom.Min
	maxZoom := source.Zoom.Max
	// if t.Zoom > maxZoom && source.HasRegion() { // only if tile max zoom always > region max zoom
	if source.HasRegion() {
		if source.Region.Polygons.Contains(t.ToLangLong()) {
			h.logger.Printf("[DEBUG] Point(%v) inside Region(%v)", t, source.Region.File)
			minZoom = source.Region.Zoom.Min
			maxZoom = source.Region.Zoom.Max
		}
	}

	if t.Zoom < minZoom || t.Zoom > maxZoom {
		h.logger.Printf("[ERROR] Wrong zoom level for Source(%v): Zoom(%v)", source.Name, t.Zoom)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	mt, err := t.Mimetype()
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	found, mtime := h.cache.Check(t, source.CacheDir)
	// found in cache
	if found {
		etag := `"` + util.DigestString(mtime.String()) + `"`
		h.replyFromCache(w, t, source.CacheDir, mt, etag, r.Header.Get("If-None-Match"))
		return
	}

	// not found in cache
	if h.cfg.Service.UseSource {
		h.replyFromSource(w, t, source.URL, mt)
	}

	// fetch tiles for metatile and write to cache?
	if h.cfg.Service.UseWriter {
		m := t.ToMetatile()
		qkey := style + "/" + m.Path()

		if h.queue.Add(qkey) {
			h.logger.Printf("[DEBUG] Add to queue: %v", qkey)
			go h.fetchAndWrite(m, style, t.Ext, source.CacheDir, source.URL, qkey)
			return
		}

		h.logger.Printf("[DEBUG] Already in queue, skip: %v", qkey)
		return
	}

	return
}

func (h mapsHandler) replyFromCache(w http.ResponseWriter, t coords.Tile, dir, mt, etag, ifNoneMatch string) {
	w.Header().Set("Etag", etag)
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%v", h.cfg.Service.MaxAge))

	if ifNoneMatch == etag {
		h.logger.Printf("[DEBUG] File not modified: Etag(%v) == If-None-Match(%v)", etag, ifNoneMatch)
		w.WriteHeader(http.StatusNotModified)
		return
	}

	data, err := h.cache.Read(t, dir)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", mt)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
	return
}

func (h mapsHandler) replyFromSource(w http.ResponseWriter, t coords.Tile, sURL, mt string) {
	url := strings.Replace(sURL, "{tile}", t.Path(), 1)
	h.logger.Printf("Get from source URL(%v)", url)

	data, err := httpclient.Get(url, h.cfg.HTTPClient.UserAgent)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", mt)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
	return
}

func (h mapsHandler) fetchAndWrite(m coords.Metatile, style, ext, cacheDir, sURL, qkey string) error {
	defer func() {
		h.logger.Printf("Done, del from queue: %v", qkey)
		h.queue.Del(qkey)
	}()

	xybox := m.ToXYBox()
	h.logger.Printf("Fetch Style(%v) Zoom(%v) X(%v-%v) Y(%v-%v) Source(%v)",
		style, m.Zoom, xybox.X[0], xybox.X[len(xybox.X)-1], xybox.Y[0], xybox.Y[len(xybox.Y)-1], sURL)

	var data [][]byte
	for _, x := range xybox.X {
		for _, y := range xybox.Y {
			tile := strconv.Itoa(m.Zoom) + "/" + strconv.Itoa(x) + "/" + strconv.Itoa(y) + `.` + ext
			url := strings.Replace(sURL, "{tile}", tile, 1)
			// fc.logger.Printf("[DEBUG] Filecache/fetchAndWrite: Fetch %v", url)
			res, err := httpclient.Get(url, h.cfg.HTTPClient.UserAgent)
			if err != nil {
				h.logger.Printf("[ERROR] fetchAndWrite: %v", err)
				return fmt.Errorf("fetchAndWrite: %v", err)
			}
			data = append(data, res)
		}
	}

	err := h.cache.Write(m, cacheDir, data)
	if err != nil {
		h.logger.Printf("[ERROR] fetchAndWrite: %v", err)
		return fmt.Errorf("fetchAndWrite: %v", err)
	}

	return nil
}
