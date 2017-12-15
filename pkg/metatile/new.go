package metatile

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/tierpod/metatiles-cacher/pkg/tile"
)

var reMetatile = regexp.MustCompile(`(\w+)/(\d+)/(\d+)/(\d+)/(\d+)/(\d+)/(\d+)\.meta`)

// NewFromURL creates Metatile from url.
func NewFromURL(url string) (Metatile, error) {
	items := reMetatile.FindStringSubmatch(url)
	if len(items) == 0 {
		return Metatile{}, fmt.Errorf("could not parse url string to Metatile struct")
	}

	zoom, _ := strconv.Atoi(items[2])
	h4, _ := strconv.Atoi(items[3])
	h3, _ := strconv.Atoi(items[4])
	h2, _ := strconv.Atoi(items[5])
	h1, _ := strconv.Atoi(items[6])
	h0, _ := strconv.Atoi(items[7])
	h := hashes{h0, h1, h2, h3, h4}

	x, y := h.XY()

	return Metatile{
		Map:    items[1],
		Zoom:   zoom,
		Hashes: h,
		X:      x,
		Y:      y,
	}, nil
}

// NewFromTile creates Metatile from Tile
func NewFromTile(t tile.Tile) Metatile {
	h := xyToHashes(t.X, t.Y)
	x, y := h.XY()
	return Metatile{
		Map:    t.Map,
		Zoom:   t.Zoom,
		Hashes: h,
		X:      x,
		Y:      y,
	}
}

func xyToHashes(x, y int) hashes {
	var xx, yy, mask int

	mask = MaxSize - 1
	xx = x & ^mask
	yy = y & ^mask
	h := hashes{}

	for i := 0; i < 5; i++ {
		h[i] = ((xx & 0x0f) << 4) | (yy & 0x0f)
		xx >>= 4
		yy >>= 4
	}

	return h
}
