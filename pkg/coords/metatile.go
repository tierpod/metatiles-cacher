package coords

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tierpod/metatiles-cacher/pkg/util"
)

// MaxMetatileSize is the maximum metatile size. Usually, metatile contains 8 * 8 tiles.
const MaxMetatileSize int = 8

// XYBox is the box of x and y coordinates contains in the metatile.
type XYBox struct {
	X []int
	Y []int
}
type hashes [5]int

// Metatile describes metatile coordinates: Zoom level and Hashes, calculated from Tile.
type Metatile struct {
	Zoom   int
	Hashes hashes
}

func (m Metatile) String() string {
	return fmt.Sprintf("Metatile{Zoom:%v Hashes:%v}", m.Zoom, m.Hashes)
}

// Size returns metatile size for this zoom level.
func (m Metatile) Size() int {
	size := metatileSize(m.Zoom)
	return size
}

// ToXYBox returns box of x, y coordinates for this metatile.
func (m Metatile) ToXYBox() XYBox {
	size := m.Size()
	xMin, yMin := metatileToXY(m.Hashes)
	x := util.MakeIntSlice(xMin, xMin+size)
	y := util.MakeIntSlice(yMin, yMin+size)
	return XYBox{x, y}
}

// MinXY returns mininal x and y coordinates for this metatile.
func (m Metatile) MinXY() (x int, y int) {
	x, y = metatileToXY(m.Hashes)
	return
}

// Path returns filepath of metatile, based on zoom level and hashes. Delimiter is "/".
func (m Metatile) Path() string {
	h0 := strconv.Itoa(m.Hashes[0])
	h1 := strconv.Itoa(m.Hashes[1])
	h2 := strconv.Itoa(m.Hashes[2])
	h3 := strconv.Itoa(m.Hashes[3])
	h4 := strconv.Itoa(m.Hashes[4])
	return strconv.Itoa(m.Zoom) + "/" + h4 + "/" + h3 + "/" + h2 + "/" + h1 + "/" + h0 + ".meta"
}

// Path returns filepath of metatile
// slow version with filepath.Join()
/*func (m Metatile) Path() string {
	p := []string{
		strconv.Itoa(m.Z),
		strconv.Itoa(m.Hashes[4]),
		strconv.Itoa(m.Hashes[3]),
		strconv.Itoa(m.Hashes[2]),
		strconv.Itoa(m.Hashes[1]),
		strconv.Itoa(m.Hashes[0]) + ".meta",
	}
	dirs = append(dirs, p...)
	return filepath.Join(dirs...)
}*/

func metatileToXY(h hashes) (x, y int) {
	var xx, yy int

	for i := 4; i >= 0; i-- {
		xx <<= 4
		yy <<= 4
		xx = xx | (h[i]&0xf0)>>4
		yy = yy | (h[i] & 0x0f)
	}

	return xx, yy
}

func metatileSize(z int) int {
	n := int(uint(1) << uint(z))
	if n < MaxMetatileSize {
		return n
	}
	return MaxMetatileSize
}

// MetatileMinURLPathItems is the minimum url items length, splitted by separator "/".
//
// Example: /style/zoom/0/0/0/0/0.meta has length 7.
const MetatileMinURLPathItems int = 7

// NewMetatileFromURL extracts Metatile zoom, hashes, style from url string.
func NewMetatileFromURL(url string) (m Metatile, style string, err error) {
	items := strings.Split(url, "/")
	il := len(items)
	if il < MetatileMinURLPathItems {
		err = fmt.Errorf("NewMetatileFromURL: wrong url items length: %v (%v/%v)", url, il, MetatileMinURLPathItems)
		return
	}

	// processing -1 value (item.format)
	last := strings.Split(items[il-1], ".")
	if len(last) != 2 {
		err = fmt.Errorf("NewMetatileFromURL: wrong last item: %v", items[il-1])
		return
	}

	h0, err := strconv.Atoi(last[0])
	if err != nil {
		err = fmt.Errorf("NewMetatileFromURL: h0: %v", err)
		return
	}

	h1, err := strconv.Atoi(items[il-2])
	if err != nil {
		err = fmt.Errorf("NewMetatileFromURL: h1: %v", err)
		return
	}

	h2, err := strconv.Atoi(items[il-3])
	if err != nil {
		err = fmt.Errorf("NewMetatileFromURL: h2: %v", err)
		return
	}

	h3, err := strconv.Atoi(items[il-4])
	if err != nil {
		err = fmt.Errorf("NewMetatileFromURL: h3: %v", err)
		return
	}

	h4, err := strconv.Atoi(items[il-5])
	if err != nil {
		err = fmt.Errorf("NewMetatileFromURL: h4: %v", err)
		return
	}
	m.Hashes = [5]int{h0, h1, h2, h3, h4}

	m.Zoom, err = strconv.Atoi(items[il-6])
	if err != nil {
		err = fmt.Errorf("NewMetatileFromURL: Zoom: %v", err)
		return
	}

	style = items[il-7]

	return m, style, nil
}
