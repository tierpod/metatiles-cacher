package fetch

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/httpclient"
	"github.com/tierpod/metatiles-cacher/pkg/metatile"
)

// ErrInQueue contains error message if item already preseint in queue.
var ErrInQueue = errors.New("item already present in queue")

func (f *Fetch) metatile(mt metatile.Metatile, URLTmpl string) (metatile.Data, error) {
	var data metatile.Data
	xybox := mt.XYBox()
	for _, x := range xybox.X {
		for _, y := range xybox.Y {
			offset := metatile.XYOffset(x, y)
			url := strings.Replace(URLTmpl, "{z}", strconv.Itoa(mt.Zoom), 1)
			url = strings.Replace(url, "{x}", strconv.Itoa(x), 1)
			url = strings.Replace(url, "{y}", strconv.Itoa(y), 1)

			res, err := httpclient.Get(url, f.cfg.UserAgent)
			if err != nil {
				return data, err
			}
			data[offset] = res
		}
	}
	// debug slow connections
	time.Sleep(time.Second * 10)
	return data, nil
}

// Metatile fetchs metatile data, using URLTmpl as URL template with placeholders {x} {y} {z}.
// If metatile already in the fetching queue, return skipped = true and empty data.
func (f *Fetch) Metatile(mt metatile.Metatile, URLTmpl string) (data metatile.Data, err error) {
	key := mt.Filepath("")

	if f.queue.HasKey(key) {
		f.logger.Printf("[DEBUG] already in queue, skip: %v", key)
		return metatile.Data{}, ErrInQueue
	}

	f.logger.Printf("[DEBUG] add to queue: %v", key)
	f.queue.Add(key)
	defer func() {
		f.logger.Printf("done, del from queue: %v", key)
		f.queue.Del(key)
	}()

	data, err = f.metatile(mt, URLTmpl)
	if err != nil {
		return metatile.Data{}, err
	}

	return data, nil
}

// MetatileWriteToCache fetchs metatile and writes it to cache. If metatile already in the fetching
// queue, wait for fetching and writing complete. Return error if timeout appears.
func (f *Fetch) MetatileWriteToCache(mt metatile.Metatile, URLTmpl string, w cache.Writer) error {
	timeout := 30 // TODO: timeout to cfg.httpclient
	key := mt.Filepath("")

	if f.queue.HasKey(key) {
		f.logger.Printf("[DEBUG] already in queue, wait: %v", key)
		if errw := f.queue.Wait(key, timeout); errw != nil {
			return errw
		}
		return nil
	}

	f.logger.Printf("[DEBUG] add to queue: %v", key)
	f.queue.Add(key)
	defer func() {
		f.logger.Printf("done, del from queue: %v", key)
		f.queue.Del(key)
	}()

	data, err := f.metatile(mt, URLTmpl)
	if err != nil {
		return err
	}

	err = w.Write(mt, data)
	if err != nil {
		return err
	}

	return nil
}
