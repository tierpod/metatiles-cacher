package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/httpclient"
	"github.com/tierpod/metatiles-cacher/pkg/metatile"
	"github.com/tierpod/metatiles-cacher/pkg/queue"
	"github.com/tierpod/metatiles-cacher/pkg/tile"
	"github.com/tierpod/metatiles-cacher/pkg/util"
)

type fetchHandler struct {
	logger *log.Logger
	cache  cache.ReadWriter
	cfg    *config.Config
	queue  *queue.Uniq
}

func (h fetchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	if t.Zoom < source.Zoom.Min || t.Zoom > source.Zoom.Max {
		h.logger.Printf("[ERROR] wrong zoom level for Source(%v): Zoom(%v)", source.Name, t.Zoom)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	_, err = util.Mimetype(t.Ext)
	if err != nil {
		h.logger.Printf("[ERROR] %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// fetch tiles for metatile and write to cache
	mt := metatile.NewFromTile(t)
	key := mt.Filepath("")

	if !h.queue.Add(key) {
		// TODO: wait to end fetching
		h.logger.Printf("[DEBUG] already in queue, skip: %v", key)
		w.WriteHeader(http.StatusCreated)
		return
	}

	h.logger.Printf("[DEBUG] add to queue: %v", key)
	err = h.fetchAndWrite(mt, t.Ext, source.URL, key)
	if err != nil {
		h.logger.Printf("[ERROR]: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (h fetchHandler) fetchAndWrite(mt metatile.Metatile, ext, sURL, key string) error {
	defer func() {
		h.logger.Printf("done, del from queue: %v", key)
		h.queue.Del(key)
	}()

	data, err := httpclient.FetchMetatile(mt, ext, sURL, h.cfg.HTTPClient.UserAgent)
	if err != nil {
		return err
	}

	err = h.cache.Write(mt, data)
	if err != nil {
		return fmt.Errorf("fetchAndWrite: %v", err)
	}

	return nil
}
