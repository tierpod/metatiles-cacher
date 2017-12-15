package httpclient

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/metatile"
	"github.com/tierpod/metatiles-cacher/pkg/queue"
	"github.com/tierpod/metatiles-cacher/pkg/tile"
)

// Fetch is the basic struct for fetcher.
type Fetch struct {
	logger *log.Logger
	queue  *queue.Uniq
	cfg    config.HTTPClient
}

// NewFetch creates new Fetch.
func NewFetch(cfg config.HTTPClient, logger *log.Logger) *Fetch {
	q := queue.NewUniq()
	return &Fetch{
		logger: logger,
		queue:  q,
		cfg:    cfg,
	}
}

// Tile fetchs tile data, using URLTmpl as URL template with placeholders {x} {y} {z}.
func (f *Fetch) Tile(t tile.Tile, URLTmpl string) (tile.Data, error) {
	url := strings.Replace(URLTmpl, "{z}", strconv.Itoa(t.Zoom), 1)
	url = strings.Replace(url, "{x}", strconv.Itoa(t.X), 1)
	url = strings.Replace(url, "{y}", strconv.Itoa(t.Y), 1)

	f.logger.Printf("Fetch/Tile: get from URL(%v)", url)

	data, err := Get(url, f.cfg.UserAgent)
	if err != nil {
		f.logger.Printf("[ERROR] Fetch/Tile: %v", err)
		return nil, fmt.Errorf("Fetch/Tile: %v", err)
	}

	return data, nil
}

// Metatile fetchs metatile data, using URLTmpl as URL template with placeholders {x} {y} {z}.
// If metatile already in the fetching queue, return skipped = true and empty data.
func (f *Fetch) Metatile(mt metatile.Metatile, URLTmpl string) (data metatile.Data, skipped bool, err error) {
	key := mt.Filepath("")

	if !f.queue.Add(key) {
		f.logger.Printf("[DEBUG] already in queue, skip: %v", key)
		return metatile.Data{}, true, nil
	}

	f.logger.Printf("[DEBUG] add to queue: %v", key)
	defer func() {
		f.logger.Printf("done, del from queue: %v", key)
		f.queue.Del(key)
	}()

	data, err = f.metatile(mt, URLTmpl)
	if err != nil {
		return metatile.Data{}, true, err
	}

	return data, false, nil
}

// MetatileWait fetchs metatile data, using URLTmpl as URL template with placeholders {x} {y} {z}.
// If metatile alreqdy in the fetching queue, wait to complete and return done = false and empty data.
func (f *Fetch) MetatileWait(mt metatile.Metatile, URLTmpl string) (data metatile.Data, done bool, err error) {
	key := mt.Filepath("")

	if f.queue.HasKey(key) {
		f.logger.Printf("[DEBUG] already in queue, skip: %v", key)
		if errw := f.queue.Wait(key, 30); errw != nil {
			return metatile.Data{}, false, errw
		}
		return metatile.Data{}, false, nil
	}

	f.queue.Add(key)
	f.logger.Printf("[DEBUG] add to queue: %v", key)
	defer func() {
		f.logger.Printf("done, del from queue: %v", key)
		f.queue.Del(key)
	}()

	data, err = f.metatile(mt, URLTmpl)
	if err != nil {
		return metatile.Data{}, false, err
	}

	return data, true, nil
}

func (f *Fetch) metatile(mt metatile.Metatile, URLTmpl string) (metatile.Data, error) {
	var data metatile.Data
	xybox := mt.XYBox()
	for _, x := range xybox.X {
		for _, y := range xybox.Y {
			offset := metatile.XYOffset(x, y)
			url := strings.Replace(URLTmpl, "{z}", strconv.Itoa(mt.Zoom), 1)
			url = strings.Replace(url, "{x}", strconv.Itoa(x), 1)
			url = strings.Replace(url, "{y}", strconv.Itoa(y), 1)

			res, err := Get(url, f.cfg.UserAgent)
			if err != nil {
				return data, err
			}
			data[offset] = res
		}
	}
	// debug slow connections
	// time.Sleep(time.Second * 10)
	return data, nil
}
