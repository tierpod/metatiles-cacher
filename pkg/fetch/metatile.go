package fetch

import (
	"errors"
	"strconv"
	"strings"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/httpclient"
	"github.com/tierpod/metatiles-cacher/pkg/metatile"
)

// ErrQueueHasKey contains error message if queue already has item with key.
var ErrQueueHasKey = errors.New("queue already has item with this key")

// Metatile fetchs metatile data, using URLTmpl as template with placeholders: {z} {x} {y}.
func (f *Fetch) Metatile(mt metatile.Metatile, URLTmpl string) (metatile.Data, error) {
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
	// time.Sleep(time.Second * 10)
	return data, nil
}

// MetatileWaitWriteToCache fetchs metatile data and writes it to cache. If metatile already in the
// fetching queue, wait for fetching and writing complete.
func (f *Fetch) MetatileWaitWriteToCache(mt metatile.Metatile, URLTmpl string, w cache.Writer) error {
	key := mt.Filepath("")

	if f.queue.HasKey(key) {
		f.logger.Printf("[DEBUG] already in queue, wait: %v", key)
		if errw := f.queue.Wait(key, f.cfg.QueueTimeout); errw != nil {
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

	data, err := f.Metatile(mt, URLTmpl)
	if err != nil {
		return err
	}

	err = w.Write(mt, data)
	if err != nil {
		return err
	}

	return nil
}

// MetatileWriteToCache fetchs metatile data and writes it to cache. If metatile already in the
// fetching queue, return error ErrQueueHasKey.
func (f *Fetch) MetatileWriteToCache(mt metatile.Metatile, URLTmpl string, w cache.Writer) error {
	key := mt.Filepath("")

	if f.queue.HasKey(key) {
		f.logger.Printf("[DEBUG] already in queue, wait: %v", key)
		return ErrQueueHasKey
	}

	f.logger.Printf("[DEBUG] add to queue: %v", key)
	f.queue.Add(key)
	defer func() {
		f.logger.Printf("done, del from queue: %v", key)
		f.queue.Del(key)
	}()

	data, err := f.Metatile(mt, URLTmpl)
	if err != nil {
		return err
	}

	err = w.Write(mt, data)
	if err != nil {
		return err
	}

	return nil
}
