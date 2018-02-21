// Package metatile provides Metatile struct with coordinates and metatile file encoder/decoder.
//
// Metatile format description:
// https://github.com/openstreetmap/mod_tile/blob/master/src/metatile.cpp
package metatile

import (
	"fmt"
	"path"
	"regexp"
	"strconv"

	"github.com/tierpod/go-osm/tile"
)

const (
	// MaxSize is the maximum metatile size.
	MaxSize int = 8
	// Area is the area of metatile.
	Area int = MaxSize * MaxSize
	// Ext is the metatile file extension.
	Ext string = ".meta"
)

// Metatile describes metatile coordinates.
type Metatile struct {
	Zoom   int
	Style  string
	X, Y   int
	hashes hashes
}

func (m Metatile) String() string {
	s := m.Size() - 1
	return fmt.Sprintf("Metatile{Zoom:%v X:%v-%v Y:%v-%v Style:%v Ext:%v}", m.Zoom, m.X, m.X+s, m.Y, m.Y+s, m.Style, Ext)
}

// Filepath returns metatile file path, based on basedir and coordinates.
func (m Metatile) Filepath(basedir string) string {
	zoom := strconv.Itoa(m.Zoom)
	h0 := strconv.Itoa(m.hashes[0]) + Ext
	h1 := strconv.Itoa(m.hashes[1])
	h2 := strconv.Itoa(m.hashes[2])
	h3 := strconv.Itoa(m.hashes[3])
	h4 := strconv.Itoa(m.hashes[4])
	return path.Join(basedir, m.Style, zoom, h4, h3, h2, h1, h0)
}

// Size return metatile size for current zoom level.
func (m Metatile) Size() int {
	n := int(uint(1) << uint(m.Zoom))
	if n < MaxSize {
		return n
	}
	return MaxSize
}

// XYBox returns arrays of x and y coordinates contains inside metatile.
func (m Metatile) XYBox() (xx []int, yy []int) {
	for x := m.X; x < m.X+m.Size(); x++ {
		xx = append(xx, x)
	}

	for y := m.Y; y < m.Y+m.Size(); y++ {
		yy = append(yy, y)
	}

	return xx, yy
}

// New creates Metatile from (z, x, y) coordinates and style.
func New(z, x, y int, style string) Metatile {
	h := xyToHashes(x, y)
	mx, my := h.xy()
	return Metatile{
		Style:  style,
		Zoom:   z,
		hashes: h,
		X:      mx,
		Y:      my,
	}
}

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

	x, y := h.xy()

	return Metatile{
		Style:  items[1],
		Zoom:   zoom,
		hashes: h,
		X:      x,
		Y:      y,
	}, nil
}

// NewFromTile creates Metatile from Tile.
func NewFromTile(t tile.Tile) Metatile {
	return New(t.Zoom, t.X, t.Y, t.Style)
}
