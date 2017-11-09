// Package cache provides interfaces for read tile from cache and write metatile to cache.
package cache

import (
	"time"

	"github.com/tierpod/metatiles-cacher/pkg/coords"
)

// Reader provides interface for read tile data from metatiles cache.
type Reader interface {
	Read(t coords.Tile, style string) (data []byte, err error)
	Check(t coords.Tile, style string) (found bool, mtime time.Time)
}

// Writer provides interface for write metatile data data to cache.
type Writer interface {
	Write(m coords.Metatile, style string, data [][]byte) error
}

// ReadWriter includes Reader and Writer interfaces.
type ReadWriter interface {
	Reader
	Writer
}
