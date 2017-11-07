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

type addHandler struct {
	logger *log.Logger
	queue  *queue.Uniq
	cache  cache.Writer
	cfg    *config.Service
}

func (h addHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	meta, style, err := coords.NewMetaFromURL(r.URL.Path)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if meta.Z < h.cfg.Zoom.Min || meta.Z > h.cfg.Zoom.Max {
		h.logger.Printf("[ERROR] Wrong zoom level: Z(%v)", meta.Z)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	source, found := h.cfg.Sources.Map[style]
	if !found {
		h.logger.Printf("[ERROR] Style not found in sources: %v", style)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	q := style + "/" + meta.Path()
	if h.queue.Add(q) {
		h.logger.Printf("[DEBUG] Add to queue: %v", q)
		go h.fetchAndWrite(meta, style, source)
		return
	}

	h.logger.Printf("[DEBUG] Already in queue, skip: %v", q)
	return
}

func (h addHandler) fetchAndWrite(meta coords.Metatile, style, source string) error {
	var result [][]byte
	var url, zxy string

	defer func() {
		q := style + "/" + meta.Path()
		h.logger.Printf("Done, del from queue: %v", q)
		h.queue.Del(q)
	}()

	minX, minY := meta.MinXY()
	h.logger.Printf("Fetch Style(%v) Z(%v) X(%v-%v) Y(%v-%v) Source(%v)",
		style, meta.Z, minX, minX+meta.Size(), minY, minY+meta.Size(), source)

	xybox := meta.ConvertToXYBox()
	for _, x := range xybox.X {
		for _, y := range xybox.Y {
			zxy = strconv.Itoa(meta.Z) + "/" + strconv.Itoa(x) + "/" + strconv.Itoa(y) + ".png"
			url = strings.Replace(source, "{zxy}", zxy, 1)
			// fc.logger.Printf("[DEBUG] Filecache/fetchAndWrite: Fetch %v", url)
			res, err := httpclient.Get(url, h.cfg.HTTPClient.UserAgent)
			if err != nil {
				h.logger.Printf("[ERROR] Add/fetchAndWrite: %v", err)
				return fmt.Errorf("Add/fetchAndWrite: %v", err)
			}
			result = append(result, res)
		}
	}

	err := h.cache.Write(meta, style, result)
	if err != nil {
		h.logger.Printf("[ERROR] Add/fetchAndWrite: %v", err)
		return fmt.Errorf("Add/fetchAndWrite: %v", err)
	}

	return nil
}
