// Package cache contains metatile cache implementation.
package cache

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/tierpod/go-osm/metatile"
	"github.com/tierpod/go-osm/tile"
	"github.com/tierpod/metatiles-cacher/pkg/config"
)

// MetatileCache is the file cache using metatile file format.
type MetatileCache struct {
	cfg    config.Cache
	logger *log.Logger
}

// NewMetatileCache creates new MetatileCache with `cfg` as configuration and `logger`.
func NewMetatileCache(cfg config.Cache, logger *log.Logger) (*MetatileCache, error) {
	if _, err := os.Stat(cfg.Dir); os.IsNotExist(err) {
		return nil, err
	}

	c := MetatileCache{
		cfg:    cfg,
		logger: logger,
	}

	return &c, nil
}

// Read reads metatile with tile `t`, decodes and return tile data.
func (c *MetatileCache) Read(t tile.Tile) ([]byte, error) {
	mt := metatile.NewFromTile(t)
	p := mt.Filepath(c.cfg.Dir)
	c.logger.Printf("[DEBUG] (cache) read %v from metatile %v", t, p)

	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder, err := metatile.NewDecoder(f)
	if err != nil {
		return nil, err
	}

	data, err := decoder.Tile(t.X, t.Y)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Write writes `data` to metatile file.
func (c *MetatileCache) Write(mt metatile.Metatile, data [][]byte) error {
	p := mt.Filepath(c.cfg.Dir)

	err := os.MkdirAll(filepath.Dir(p), 0777)
	if err != nil {
		return err
	}

	tmpDir := filepath.Dir(p)
	f, err := ioutil.TempFile(tmpDir, "write")
	if err != nil {
		return err
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	c.logger.Printf("[DEBUG] (cache) write to temp file: %v", f.Name())
	err = mt.EncodeWrite(f, data)
	if err != nil {
		return err
	}

	err = os.Rename(f.Name(), p)
	if err != nil {
		return err
	}

	if err := os.Chmod(p, 0666); err != nil {
		return err
	}

	c.logger.Printf("[INFO] (cache) wrote %v", p)
	return nil
}

// Check checks if metatile file for tile `t` exists on disk. Returns modification time and
// found=true if exists, otherwise found=false.
func (c *MetatileCache) Check(t tile.Tile) (mtime time.Time, found bool) {
	mt := metatile.NewFromTile(t)
	p := mt.Filepath(c.cfg.Dir)
	c.logger.Printf("[DEBUG] (cache) check %v", p)

	stat, err := os.Stat(p)
	if !os.IsNotExist(err) {
		return stat.ModTime(), true
	}

	return time.Time{}, false
}
