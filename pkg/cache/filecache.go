package cache

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/metatile"
	"github.com/tierpod/metatiles-cacher/pkg/tile"
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
func (fc *FileCache) Read(t tile.Tile) (data tile.Data, err error) {
	mt := metatile.NewFromTile(t)
	path := mt.Filepath(fc.cfg.RootDir)
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
func (fc *FileCache) Check(t tile.Tile) (found bool, mtime time.Time) {
	mt := metatile.NewFromTile(t)
	path := mt.Filepath(fc.cfg.RootDir)
	fc.logger.Printf("[DEBUG] FileCache: check %v", path)

	stat, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return true, stat.ModTime()
	}

	return false, time.Time{}
}

// Write writes metatile data to disk.
func (fc *FileCache) Write(m metatile.Metatile, data metatile.Data) error {
	path := m.Filepath(fc.cfg.RootDir)
	fc.logger.Printf("FileCache: write %v", path)

	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return fmt.Errorf("FileCache: %v", err)
	}

	tmpDir := filepath.Dir(path)
	f, err := ioutil.TempFile(tmpDir, "write")
	if err != nil {
		return fmt.Errorf("FileCache: %v", err)
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	// fc.logger.Printf("[DEBUG] FileCache: write to temp file: %v", file.Name())

	/*err = metatile.Encode(file, m, data)
	if err != nil {
		return fmt.Errorf("FileCache: %v", err)
	}*/
	err = m.Encode(f, data)
	if err != nil {
		return fmt.Errorf("FileCache: %v", err)
	}

	err = os.Rename(f.Name(), path)
	if err != nil {
		return fmt.Errorf("FileCache: %v", err)
	}

	if err := os.Chmod(path, 0666); err != nil {
		return fmt.Errorf("FileCache: %v", err)
	}

	return nil
}
