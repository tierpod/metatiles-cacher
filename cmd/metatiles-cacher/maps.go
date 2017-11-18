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

	source, err := h.cfg.Source(style)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if t.Zoom < source.Zoom.Min || t.Zoom > source.Zoom.Max {
		h.logger.Printf("[ERROR] Wrong zoom level for Source(%v): Zoom(%v)", source.Name, t.Zoom)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if source.HasRegion() {
		h.logger.Println("[DEBUG] Check if tile coords in given region")
		for _, p := range source.Region.Polygons {
			in := p.Contains(t.ToLangLong())
			if !in {
				h.logger.Printf("[WARN] Point not in given region")
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}
	}

	_, err = t.Mimetype()
	if err != nil {
		h.logger.Printf("[ERROR] Wrong extension: Ext(%v)", t.Ext)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.logger.Printf("[DEBUG] Got request %v style(%v)", t, style)

	found, mtime := h.cache.Check(t, source.CacheDir)
	// found in cache
	if found {
		etag := `"` + util.DigestString(mtime.String()) + `"`
		h.replyFromCache(w, t, source.CacheDir, etag, r.Header.Get("If-None-Match"))
		return
	}

	// not found in cache
	if h.cfg.Service.UseSource {
		h.replyFromSource(w, t, source.URL)
	}

	// fetch tiles for metatile and write to cache?
	if h.cfg.Service.UseWriter {
		m := t.ToMetatile()
		qkey := style + "/" + m.Path()

		if h.queue.Add(qkey) {
			h.logger.Printf("[DEBUG] Add to queue: %v", qkey)
			go h.fetchAndWrite(m, t.Ext, source.CacheDir, source.URL, qkey)
			return
		}

		h.logger.Printf("[DEBUG] Already in queue, skip: %v", qkey)
		return
	}

	return
}

func (h mapsHandler) replyFromCache(w http.ResponseWriter, t coords.Tile, dir string, etag, ifNoneMatch string) {
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

	mt, _ := t.Mimetype()
	w.Header().Set("Content-Type", mt)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
	return
}

func (h mapsHandler) replyFromSource(w http.ResponseWriter, t coords.Tile, source string) {
	url := strings.Replace(source, "{tile}", t.Path(), 1)
	h.logger.Printf("Get from Source(%v)", url)

	data, err := httpclient.Get(url, h.cfg.HTTPClient.UserAgent)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mt, _ := t.Mimetype()
	w.Header().Set("Content-Type", mt)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
	return
}

func (h mapsHandler) fetchAndWrite(m coords.Metatile, ext, cacheDir, sourceURL, qkey string) error {
	var result [][]byte
	var url string

	defer func() {
		h.logger.Printf("Done, del from queue: %v", qkey)
		h.queue.Del(qkey)
	}()

	minX, minY := m.MinXY()
	h.logger.Printf("Fetch Style(%v) Zoom(%v) X(%v-%v) Y(%v-%v) Source(%v)",
		cacheDir, m.Zoom, minX, minX+m.Size(), minY, minY+m.Size(), sourceURL)

	xybox := m.ToXYBox()
	for _, x := range xybox.X {
		for _, y := range xybox.Y {
			tile := strconv.Itoa(m.Zoom) + "/" + strconv.Itoa(x) + "/" + strconv.Itoa(y) + `.` + ext
			url = strings.Replace(sourceURL, "{tile}", tile, 1)
			// fc.logger.Printf("[DEBUG] Filecache/fetchAndWrite: Fetch %v", url)
			res, err := httpclient.Get(url, h.cfg.HTTPClient.UserAgent)
			if err != nil {
				h.logger.Printf("[ERROR] fetchAndWrite: %v", err)
				return fmt.Errorf("fetchAndWrite: %v", err)
			}
			result = append(result, res)
		}
	}

	err := h.cache.Write(m, cacheDir, result)
	if err != nil {
		h.logger.Printf("[ERROR] fetchAndWrite: %v", err)
		return fmt.Errorf("fetchAndWrite: %v", err)
	}

	return nil
}
