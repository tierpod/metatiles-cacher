package cache

import (
	"fmt"
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
	cfg    config.FileCacheSection
	logger *log.Logger
}

// NewFileCache creates new FileCache. Return error if cfg.RootDir does not exists.
func NewFileCache(cfg config.FileCacheSection, logger *log.Logger) (*FileCache, error) {
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
func (fc *FileCache) Read(t coords.Tile, style string) (data []byte, err error) {
	path := fc.cfg.RootDir + "/" + style + "/" + t.ToMetatile().Path()
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
func (fc *FileCache) Check(t coords.Tile, style string) (found bool, mtime time.Time) {
	path := fc.cfg.RootDir + "/" + style + "/" + t.ToMetatile().Path()
	fc.logger.Printf("[DEBUG] FileCache: check %v", path)

	stat, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return true, stat.ModTime()
	}

	return false, time.Time{}
}

// Write writes metatile data to disk.
func (fc *FileCache) Write(m coords.Metatile, style string, data [][]byte) error {
	path := fc.cfg.RootDir + "/" + style + "/" + m.Path()
	fc.logger.Printf("FileCache: write %v", path)

	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return fmt.Errorf("FileCache: %v", err)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("FileCache: %v", err)
	}
	defer file.Close()

	err = metatile.Encode(file, m, data)
	if err != nil {
		return fmt.Errorf("FileCache: %v", err)
	}

	return nil
}
