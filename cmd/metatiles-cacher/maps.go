package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/httpclient"
	"github.com/tierpod/metatiles-cacher/pkg/latlong"
	"github.com/tierpod/metatiles-cacher/pkg/metatile"
	"github.com/tierpod/metatiles-cacher/pkg/queue"
	"github.com/tierpod/metatiles-cacher/pkg/tile"
	"github.com/tierpod/metatiles-cacher/pkg/util"
)

type mapsHandler struct {
	logger *log.Logger
	cache  cache.ReadWriter
	cfg    *config.Config

	queue *queue.Uniq
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

	mt, err := util.Mimetype(t.Ext)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	found, mtime := h.cache.Check(t)
	// found in cache
	if found {
		etag := `"` + util.DigestString(mtime.String()) + `"`
		h.replyFromCache(w, t, mt, etag, r.Header.Get("If-None-Match"))
		return
	}

	// not found in cache
	if h.cfg.Service.UseSource {
		h.replyFromSource(w, t, source.URL, mt)
	}

	// fetch tiles for metatile and write to cache?
	if h.cfg.Service.UseWriter {
		mt := metatile.NewFromTile(t)
		key := mt.Filepath("")

		if h.queue.Add(key) {
			h.logger.Printf("[DEBUG] add to queue: %v", key)
			go h.fetchAndWrite(mt, t.Ext, source.URL, key)
			return
		}

		h.logger.Printf("[DEBUG] already in queue, skip: %v", key)
		return
	}

	return
}

func (h mapsHandler) replyFromCache(w http.ResponseWriter, t tile.Tile, mt, etag, ifNoneMatch string) {
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

	w.Header().Set("Content-Type", mt)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
	return
}

func (h mapsHandler) replyFromSource(w http.ResponseWriter, t tile.Tile, sURL, mt string) {
	tile := fmt.Sprintf("%v/%v/%v%v", t.Zoom, t.X, t.Y, t.Ext)
	url := strings.Replace(sURL, "{tile}", tile, 1)
	h.logger.Printf("replyFromSource: get from URL(%v)", url)

	data, err := httpclient.Get(url, h.cfg.HTTPClient.UserAgent)
	if err != nil {
		h.logger.Printf("[ERROR] replyFromSource: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", mt)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
	return
}

func (h mapsHandler) fetchAndWrite(mt metatile.Metatile, ext, sURL, qkey string) error {
	defer func() {
		h.logger.Printf("fetchAndWrite: done, del from queue: %v", qkey)
		h.queue.Del(qkey)
	}()

	xybox := mt.XYBox()
	h.logger.Printf("fetchAndWrite: fetch Map(%v) Zoom(%v) X(%v-%v) Y(%v-%v) Source(%v)",
		mt.Map, mt.Zoom, xybox.X[0], xybox.X[len(xybox.X)-1], xybox.Y[0], xybox.Y[len(xybox.Y)-1], sURL)

	data, err := httpclient.FetchMetatile(mt, ext, sURL, h.cfg.HTTPClient.UserAgent)
	if err != nil {
		return err
	}

	err = h.cache.Write(mt, data)
	if err != nil {
		return fmt.Errorf("[ERROR] fetchAndWrite: %v", err)
	}

	return nil
}
