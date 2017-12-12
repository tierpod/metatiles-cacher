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
)

type fetchHandler struct {
	logger *log.Logger
	cache  cache.ReadWriter
	cfg    *config.Config

	queue *queue.Uniq
}

func (h fetchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	if t.Zoom < source.Zoom.Min || t.Zoom > source.Zoom.Max {
		h.logger.Printf("[ERROR] Wrong zoom level for Source(%v): Zoom(%v)", source.Name, t.Zoom)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	_, err = t.Mimetype()
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// fetch tiles for metatile and write to cache
	m := t.ToMetatile()
	qkey := style + "/" + m.Path()

	if !h.queue.Add(qkey) {
		// TODO: wait to end fetching
		h.logger.Printf("[DEBUG] Already in queue, skip: %v", qkey)
		w.WriteHeader(http.StatusCreated)
		return
	}

	h.logger.Printf("[DEBUG] Add to queue: %v", qkey)
	err = h.fetchAndWrite(m, style, t.Ext, source.CacheDir, source.URL, qkey)
	if err != nil {
		h.logger.Printf("[ERROR]: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (h fetchHandler) fetchAndWrite(m coords.Metatile, style, ext, cacheDir, sURL, qkey string) error {
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
			res, err := httpclient.Get(url, h.cfg.HTTPClient.UserAgent)
			if err != nil {
				return fmt.Errorf("fetchAndWrite: %v", err)
			}
			data = append(data, res)
		}
	}

	err := h.cache.Write(m, cacheDir, data)
	if err != nil {
		return fmt.Errorf("fetchAndWrite: %v", err)
	}

	return nil
}
