// Package bbox contains functions for make box of tiles.
package bbox

import (
	"github.com/tierpod/metatiles-cacher/pkg/latlong"
	"github.com/tierpod/metatiles-cacher/pkg/tile"
)

// NewFromLatLong calculates bbox from top to bottom coordinates for each zoom level in zooms.
// Returns output chan with tiles from this bbox.
func NewFromLatLong(zooms []int, top latlong.LatLong, bottom latlong.LatLong, ext string) <-chan (tile.Tile) {
	ch := make(chan tile.Tile)

	go func() {
		defer close(ch)
		for _, z := range zooms {
			tTop := tile.NewFromLatLong(top, z)
			tBottom := tile.NewFromLatLong(bottom, z)
			for x := tTop.X; x <= tBottom.X; x++ {
				for y := tTop.Y; y <= tBottom.Y; y++ {
					t := tile.Tile{Zoom: z, X: x, Y: y, Ext: ext}
					ch <- t
				}
			}
		}
	}()

	return ch
}
