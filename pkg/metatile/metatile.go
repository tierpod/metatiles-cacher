package metatile

import (
	"fmt"
	"path"
	"strconv"
)

// MaxSize is the maximum metatile size.
const (
	MaxSize int    = 8
	Area    int    = MaxSize * MaxSize
	Ext     string = ".meta"
)

type hashes [5]int

func (h hashes) XY() (int, int) {
	var x, y int

	for i := 4; i >= 0; i-- {
		x <<= 4
		y <<= 4
		x = x | (h[i]&0xf0)>>4
		y = y | (h[i] & 0x0f)
	}

	return x, y
}

// Metatile describes metatile coordinates: Zoom level and Hashes, calculated from Tile.
type Metatile struct {
	Zoom   int
	Map    string
	Hashes hashes
	X, Y   int
}

func (m Metatile) String() string {
	return fmt.Sprintf("Metatile{Zoom:%v Hashes:%v Map:%v Ext:%v X:%v Y:%v}", m.Zoom, m.Hashes, m.Map, Ext, m.X, m.Y)
}

// Filepath returns metatile file path, based on basedir and coordinates.
func (m Metatile) Filepath(basedir string) string {
	zoom := strconv.Itoa(m.Zoom)
	h0 := strconv.Itoa(m.Hashes[0]) + Ext
	h1 := strconv.Itoa(m.Hashes[1])
	h2 := strconv.Itoa(m.Hashes[2])
	h3 := strconv.Itoa(m.Hashes[3])
	h4 := strconv.Itoa(m.Hashes[4])
	return path.Join(basedir, m.Map, zoom, h4, h3, h2, h1, h0)
}
