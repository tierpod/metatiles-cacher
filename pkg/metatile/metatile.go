// Package metatile contains format description, provides functions for decoding and encoding
// metatile files.
//
// Metatile format description:
// https://github.com/openstreetmap/mod_tile/blob/master/src/metatile.cpp
package metatile

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strconv"

	"github.com/tierpod/metatiles-cacher/pkg/tile"
	"github.com/tierpod/metatiles-cacher/pkg/util"
)

const (
	// MaxSize is the maximum metatile size.
	MaxSize int = 8
	// Area is the area of metatile.
	Area int = MaxSize * MaxSize
	// Ext is the metatile file extension.
	Ext string = ".meta"
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

func (mt Metatile) String() string {
	return fmt.Sprintf("Metatile{Zoom:%v Hashes:%v Map:%v Ext:%v X:%v Y:%v}", mt.Zoom, mt.Hashes, mt.Map, Ext, mt.X, mt.Y)
}

// Filepath returns metatile file path, based on basedir and coordinates.
func (mt Metatile) Filepath(basedir string) string {
	zoom := strconv.Itoa(mt.Zoom)
	h0 := strconv.Itoa(mt.Hashes[0]) + Ext
	h1 := strconv.Itoa(mt.Hashes[1])
	h2 := strconv.Itoa(mt.Hashes[2])
	h3 := strconv.Itoa(mt.Hashes[3])
	h4 := strconv.Itoa(mt.Hashes[4])
	return path.Join(basedir, mt.Map, zoom, h4, h3, h2, h1, h0)
}

// Size return metatile size for current zoom level.
func (mt Metatile) Size() int {
	n := int(uint(1) << uint(mt.Zoom))
	if n < MaxSize {
		return n
	}
	return MaxSize
}

// DecodeTile reads metatile data from r, decode tile data with (x, y) coordinates and writes it to w.
func (mt Metatile) DecodeTile(w io.Writer, r io.ReadSeeker, x, y int) error {
	data, err := decodeTile(r, x, y)
	if err != nil {
		return err
	}

	io.Copy(w, bytes.NewReader(data))
	return nil
}

// EncodeTiles encodes tiles data to metatile format and writes it to w.
func (mt Metatile) EncodeTiles(w io.Writer, data Data) error {
	err := encodeMetatile(w, data, mt.X, mt.Y, mt.Zoom)
	if err != nil {
		return err
	}

	return nil
}

// XYBox is the box of (x, y) coordinates contains inside metatile.
type XYBox struct {
	X []int
	Y []int
}

// XYBox return box of (x, y) coordinates contains inside this metatile.
func (mt Metatile) XYBox() XYBox {
	size := mt.Size()
	x := util.MakeIntSlice(mt.X, mt.X+size)
	y := util.MakeIntSlice(mt.Y, mt.Y+size)
	return XYBox{x, y}
}

// Data is array of tile data.
type Data [Area]tile.Data

// XYOffset returns offset of tile data inside metatile.
func XYOffset(x, y int) int {
	mask := MaxSize - 1
	return (x&mask)*MaxSize + (y & mask)
}
