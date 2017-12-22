package cache

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/tierpod/metatiles-cacher/pkg/metatile"
	"github.com/tierpod/metatiles-cacher/pkg/tile"
)

// MetatileCacher is the interface for decoding tile data from metatile and writing metatile data to
// metatile.
type MetatileCacher interface {
	Filepath(basedir string) string
}

// MetatileCache is the filecache using metatile file format.
type MetatileCache struct {
	dir    string
	logger *log.Logger
}

// NewMetatileCache created new metatile cache with dir as root directory and logger.
func NewMetatileCache(dir string, logger *log.Logger) (*MetatileCache, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, err
	}

	c := MetatileCache{
		dir:    dir,
		logger: logger,
	}

	return &c, nil
}

// Read decodes tile data with (t.X, t.Y) coordinates from metatile file and writes it to w.
func (c *MetatileCache) Read(mc MetatileCacher, t tile.Tile, w io.Writer) error {
	p := mc.Filepath(c.dir)
	c.logger.Printf("[DEBUG] cache: read %v from metatile %v", t, p)

	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()

	var buf bytes.Buffer
	err = metatile.DecodeTileTo(&buf, f, t.X, t.Y)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, &buf)
	if err != nil {
		return err
	}
	return nil
}

// Write reads encoded metatile data from r and write it to metatile file.
func (c *MetatileCache) Write(mc MetatileCacher, r io.Reader) error {
	p := mc.Filepath(c.dir)
	c.logger.Printf("cache: write %v", p)

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

	c.logger.Printf("[DEBUG] cache: write to temp file: %v", f.Name())
	_, err = io.Copy(f, r)
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

	c.logger.Printf("[DEBUG] cache: wrote %v", p)

	return nil
}

// Check checks if metatile exists on disk. If exist, returns modification time of file and found =
// true.
func (c *MetatileCache) Check(mc MetatileCacher) (mtime time.Time, found bool) {
	p := mc.Filepath(c.dir)
	c.logger.Printf("[DEBUG] cache: check %v", p)

	stat, err := os.Stat(p)
	if !os.IsNotExist(err) {
		return stat.ModTime(), true
	}

	return time.Time{}, false
}
