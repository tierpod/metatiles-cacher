package cache

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/coords"
	"github.com/tierpod/metatiles-cacher/pkg/metatile"
)

// FileCache is the file cache struct, contains self configuration (RootDir) and logger. Implemets
// cache.ReadWriter interface.
type FileCache struct {
	cfg    config.FileCache
	logger *log.Logger
}

// NewFileCache creates new FileCache. Return error if cfg.RootDir does not exists.
func NewFileCache(cfg config.FileCache, logger *log.Logger) (*FileCache, error) {
	if _, err := os.Stat(cfg.RootDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("NewFileCache: %v is not exist", cfg.RootDir)
	}

	fc := FileCache{
		cfg:    cfg,
		logger: logger,
	}

	return &fc, nil
}

// Read reads tile data from metatile.
func (fc *FileCache) Read(t coords.Tile, dir string) (data coords.TileData, err error) {
	path := fc.cfg.RootDir + "/" + dir + "/" + t.ToMetatile().Path()
	fc.logger.Printf("[DEBUG] FileCache: read %v from metatile %v", t, path)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("FileCache: %v", err)
	}
	defer file.Close()

	data, err = metatile.GetTile(file, t)
	if err != nil {
		return nil, fmt.Errorf("FileCache: %v", err)
	}

	return data, nil
}

// Check checks if tile in the file cache. If found, return found = true and mtime = modification time of file.
func (fc *FileCache) Check(t coords.Tile, dir string) (found bool, mtime time.Time) {
	path := fc.cfg.RootDir + "/" + dir + "/" + t.ToMetatile().Path()
	fc.logger.Printf("[DEBUG] FileCache: check %v", path)

	stat, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return true, stat.ModTime()
	}

	return false, time.Time{}
}

// Write writes metatile data to disk.
func (fc *FileCache) Write(m coords.Metatile, dir string, data coords.MetatileData) error {
	path := fc.cfg.RootDir + "/" + dir + "/" + m.Path()
	fc.logger.Printf("FileCache: write %v", path)

	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return fmt.Errorf("FileCache: %v", err)
	}

	tmpDir := fc.cfg.RootDir + "/" + dir + "/" + m.Dir()
	file, err := ioutil.TempFile(tmpDir, "fetch")
	if err != nil {
		return fmt.Errorf("FileCache: %v", err)
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()
	// fc.logger.Printf("[DEBUG] FileCache: write to temp file: %v", file.Name())

	err = metatile.Encode(file, m, data)
	if err != nil {
		return fmt.Errorf("FileCache: %v", err)
	}

	err = os.Rename(file.Name(), path)
	if err != nil {
		return fmt.Errorf("FileCache: %v", err)
	}

	if err := os.Chmod(path, 0666); err != nil {
		return fmt.Errorf("FileCache: %v", err)
	}

	return nil
}
