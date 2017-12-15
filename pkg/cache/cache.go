// Package cache provides interfaces for read tile from cache and write metatile to cache.
package cache

import (
	"time"

	"github.com/tierpod/metatiles-cacher/pkg/metatile"
	"github.com/tierpod/metatiles-cacher/pkg/tile"
)

// Reader provides interface for read tile data from metatiles cache.
type Reader interface {
	Read(t tile.Tile) (data tile.Data, err error)
	Check(t tile.Tile) (found bool, mtime time.Time)
}

// Writer provides interface for write metatile data data to cache.
type Writer interface {
	Write(m metatile.Metatile, data metatile.Data) error
}

// ReadWriter includes Reader and Writer interfaces.
type ReadWriter interface {
	Reader
	Writer
}
