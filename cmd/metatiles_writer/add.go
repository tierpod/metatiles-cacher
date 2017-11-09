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
	queue *queue.Uniq
	qkey  string

	logger *log.Logger
	cache  cache.Writer
	cfg    *config.Service
}

func (h addHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m, style, err := coords.NewMetatileFromURL(r.URL.Path)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if m.Zoom < h.cfg.Zoom.Min || m.Zoom > h.cfg.Zoom.Max {
		h.logger.Printf("[ERROR] Wrong zoom level: Zoom(%v)", m.Zoom)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	source, found := h.cfg.Sources.Map[style]
	if !found {
		h.logger.Printf("[ERROR] Style not found in sources: %v", style)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.qkey = style + "/" + m.Path()
	if h.queue.Add(h.qkey) {
		h.logger.Printf("[DEBUG] Add to queue: %v", h.qkey)
		go h.fetchAndWrite(m, style, source)
		return
	}

	h.logger.Printf("[DEBUG] Already in queue, skip: %v", h.qkey)
	return
}

func (h addHandler) fetchAndWrite(m coords.Metatile, style, source string) error {
	var result [][]byte
	var url, t string

	defer func() {
		h.logger.Printf("Done, del from queue: %v", h.qkey)
		h.queue.Del(h.qkey)
	}()

	minX, minY := m.MinXY()
	h.logger.Printf("Fetch Style(%v) Zoom(%v) X(%v-%v) Y(%v-%v) Source(%v)",
		style, m.Zoom, minX, minX+m.Size(), minY, minY+m.Size(), source)

	xybox := m.ToXYBox()
	for _, x := range xybox.X {
		for _, y := range xybox.Y {
			t = strconv.Itoa(m.Zoom) + "/" + strconv.Itoa(x) + "/" + strconv.Itoa(y) + ".png"
			url = strings.Replace(source, "{tile}", t, 1)
			// fc.logger.Printf("[DEBUG] Filecache/fetchAndWrite: Fetch %v", url)
			res, err := httpclient.Get(url, h.cfg.HTTPClient.UserAgent)
			if err != nil {
				h.logger.Printf("[ERROR] Add/fetchAndWrite: %v", err)
				return fmt.Errorf("Add/fetchAndWrite: %v", err)
			}
			result = append(result, res)
		}
	}

	err := h.cache.Write(m, style, result)
	if err != nil {
		h.logger.Printf("[ERROR] Add/fetchAndWrite: %v", err)
		return fmt.Errorf("Add/fetchAndWrite: %v", err)
	}

	return nil
}
