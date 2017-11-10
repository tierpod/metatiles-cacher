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
	cfg    *config.Service

	queue *queue.Uniq
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

	// fetch tiles for metatile and write to cache?
	if h.cfg.Reader.UseWriter {
		m := t.ToMetatile()
		qkey := style + "/" + m.Path()

		if h.queue.Add(qkey) {
			h.logger.Printf("[DEBUG] Add to queue: %v", qkey)
			go h.fetchAndWrite(m, style, source, qkey)
			return
		}

		h.logger.Printf("[DEBUG] Already in queue, skip: %v", qkey)
		return
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

func (h mapsHandler) fetchAndWrite(m coords.Metatile, style, source, qkey string) error {
	var result [][]byte
	var url string

	defer func() {
		h.logger.Printf("Done, del from queue: %v", qkey)
		h.queue.Del(qkey)
	}()

	minX, minY := m.MinXY()
	h.logger.Printf("Fetch Style(%v) Zoom(%v) X(%v-%v) Y(%v-%v) Source(%v)",
		style, m.Zoom, minX, minX+m.Size(), minY, minY+m.Size(), source)

	xybox := m.ToXYBox()
	for _, x := range xybox.X {
		for _, y := range xybox.Y {
			tile := strconv.Itoa(m.Zoom) + "/" + strconv.Itoa(x) + "/" + strconv.Itoa(y) + ".png"
			url = strings.Replace(source, "{tile}", tile, 1)
			// fc.logger.Printf("[DEBUG] Filecache/fetchAndWrite: Fetch %v", url)
			res, err := httpclient.Get(url, h.cfg.HTTPClient.UserAgent)
			if err != nil {
				h.logger.Printf("[ERROR] fetchAndWrite: %v", err)
				return fmt.Errorf("fetchAndWrite: %v", err)
			}
			result = append(result, res)
		}
	}

	err := h.cache.Write(m, style, result)
	if err != nil {
		h.logger.Printf("[ERROR] fetchAndWrite: %v", err)
		return fmt.Errorf("fetchAndWrite: %v", err)
	}

	return nil
}
