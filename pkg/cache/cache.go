// Package cache provides cache interfaces and implemetation
package cache

import "github.com/tierpod/metatiles-cacher/pkg/coords"

// Reader provides interface for read data from metatiles cache
type Reader interface {
	Read(tile coords.ZXY, style string) (data []byte, found bool, err error)
}

// Writer provides interface for write data to metatiles cache
type Writer interface {
	Write(meta coords.Metatile, style string, data [][]byte) error
}

// ReadWriter includes Reader and Writer
type ReadWriter interface {
	Reader
	Writer
}
