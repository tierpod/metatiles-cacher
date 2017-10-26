package cache

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tierpod/metatiles-cacher/pkg/coords"
)

// FileCacheReader is the basic structure for file cache reader
type FileCacheReader struct {
	RootDir string
	logger  *log.Logger
}

// NewFileCacheReader creates new instance of FileCache
func NewFileCacheReader(rootDir string, logger *log.Logger) (*FileCacheReader, error) {
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("NewFileCacheReader: %v is not exist", rootDir)
	}

	fc := &FileCacheReader{
		RootDir: rootDir,
		logger:  logger,
	}

	return fc, nil
}

// Read reads tile data from metatile. Returns found = false and err = nil if metatile not found.
func (r *FileCacheReader) Read(tile coords.ZXY, style string) (data []byte, found bool, err error) {
	path := r.RootDir + "/" + style + "/" + tile.ConvertToMeta().Path()
	r.logger.Printf("[DEBUG] FileCacheReader: Get tile %v from metatile %v", tile, path)

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		r.logger.Printf("[WARN] FileCacheReader: File not found %v", path)
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("FileCacheReader: %v", err)
	}
	defer file.Close()

	data, err = GetTileFromMetatile(file, tile)
	if err != nil {
		return nil, false, fmt.Errorf("FileCacheReader: %v", err)
	}

	return data, true, nil
}

// FileCacheWriter is the basic structure for file cache writer
type FileCacheWriter struct {
	RootDir string
	logger  *log.Logger
}

// NewFileCacheWriter creates new instance of FileCache
func NewFileCacheWriter(rootDir string, logger *log.Logger) (*FileCacheWriter, error) {
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		logger.Printf("NewFileCacheWriter: Creating RootDir %v", rootDir)
	}

	err := os.MkdirAll(rootDir, 0777)
	if err != nil {
		return nil, fmt.Errorf("NewFileCacheWriter: %v", err)
	}

	fc := &FileCacheWriter{
		RootDir: rootDir,
		logger:  logger,
	}

	return fc, nil
}

// Write saves metatile data to disk
func (w *FileCacheWriter) Write(meta coords.Metatile, style string, data [][]byte) error {
	path := w.RootDir + "/" + style + "/" + meta.Path()
	w.logger.Printf("FileCacheWriter: %v", path)

	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return fmt.Errorf("FileCacheWriter: %v", err)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("FileCacheWriter: %v", err)
	}
	defer file.Close()

	err = EncodeMetatile(file, meta, data)
	if err != nil {
		return fmt.Errorf("FileCacheWriter: %v", err)
	}

	return nil
}
