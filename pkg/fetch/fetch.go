// Package fetch provides fetcher service who can fetch tile data, metatile data and writes it to
// cache.
package fetch

import (
	"log"

	"github.com/tierpod/metatiles-cacher/pkg/cache"
	"github.com/tierpod/metatiles-cacher/pkg/config"
	"github.com/tierpod/metatiles-cacher/pkg/metatile"
	"github.com/tierpod/metatiles-cacher/pkg/queue"
)

// Fetcher provides interface for fetch tile and metatile data.
type Fetcher interface {
	Metatile(mt metatile.Metatile, URLTmpl string) (data metatile.Data, err error)
}

// CacheWriter provides interface for fetch and write tile and metatile data.
type CacheWriter interface {
	Fetcher
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
