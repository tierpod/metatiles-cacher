package coords

import (
	"fmt"
	"strconv"

	"github.com/tierpod/metatiles-cacher/pkg/utils"
)

// MaxMetatileSize is the maximum metatile size. Usually, metatile contains 8 * 8 tiles.
const MaxMetatileSize int = 8

// XYBox is the box of x and y coordinates contains in the metatile.
type XYBox struct {
	X []int
	Y []int
}
type hashes [5]int

// Metatile describes metatile coordinates. Z: Zoom level, Hashes: hashes, calculated from ZXY.
type Metatile struct {
	Z      int
	Hashes hashes
}

func (m Metatile) String() string {
	return fmt.Sprintf("Metatile{Z:%v Hashes:%v}", m.Z, m.Hashes)
}

// Size returns metatile size for this zoom level.
func (m Metatile) Size() int {
	size := metatileSize(m.Z)
	return size
}

// ConvertToXYBox returns box of x, y coordinates.
func (m Metatile) ConvertToXYBox() XYBox {
	size := m.Size()
	xMin, yMin := metaToXY(m.Hashes)
	x := utils.MakeIntRange(xMin, xMin+size)
	y := utils.MakeIntRange(yMin, yMin+size)
	return XYBox{x, y}
}

// MinXY returns mininal x and y coordinates contains in metatile.
func (m Metatile) MinXY() (x int, y int) {
	x, y = metaToXY(m.Hashes)
	return
}

// Path returns filepath of metatile, based on zoom level and hashes. Delimiter is "/".
func (m Metatile) Path() string {
	h0 := strconv.Itoa(m.Hashes[0])
	h1 := strconv.Itoa(m.Hashes[1])
	h2 := strconv.Itoa(m.Hashes[2])
	h3 := strconv.Itoa(m.Hashes[3])
	h4 := strconv.Itoa(m.Hashes[4])
	return strconv.Itoa(m.Z) + "/" + h4 + "/" + h3 + "/" + h2 + "/" + h1 + "/" + h0 + ".meta"
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

func xyToMeta(x, y int) hashes {
	var xx, yy, mask int

	mask = MaxMetatileSize - 1
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

func metaToXY(h hashes) (x, y int) {
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
