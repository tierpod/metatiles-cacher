// Package fetch provides fetch service who can fetch tile/metatile data, waits for complete and
// writes it to cache.
package fetch

import (
	"log"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/metatile"
	"github.com/tierpod/metatiles-cacher/pkg/queue"
	"github.com/tierpod/metatiles-cacher/pkg/tile"
)

// Fetcher provides interface for fetch tile and metatile data.
type Fetcher interface {
	Tile(t tile.Tile, URLTmpl string) (tile.Data, error)
	Metatile(mt metatile.Metatile, URLTmpl string) (metatile.Data, error)
}

// CacheWaitWriter provides interface for fetching metatile data, writing it to cache and waiting
// for complete. All metatiles stored in fetching queue. If metatile already in queue, do not run
// new fetching, waiting for complete.
type CacheWaitWriter interface {
	// TileWaitWriteToCache(t tile.Tile, URLTmpl string, w cache.Writer) error
	MetatileWaitWriteToCache(mt metatile.Metatile, URLTmpl string, w cache.Writer) error
}

// CacheWriter provides interface for fetching metatile data and writing it to cache. All metatiles
// stored in fetching queue. If metatile already in queue, do not run new fetching, return
// ErrQueueHasKey.
type CacheWriter interface {
	// TileWriteToCache(t tile.Tile, URLTmpl string, w cache.Writer) error
	MetatileWriteToCache(mt metatile.Metatile, URLTmpl string, w cache.Writer) error
}

// Fetch is the basic struct for fetcher.
type Fetch struct {
	logger *log.Logger
	queue  *queue.Uniq
	cfg    config.HTTPClient
}

// New creates new Fetch.
func New(cfg config.HTTPClient, logger *log.Logger) *Fetch {
	q := queue.NewUniq()
	return &Fetch{
		logger: logger,
		queue:  q,
		cfg:    cfg,
	}
}
