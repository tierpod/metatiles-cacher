package coords

import (
	"fmt"
	"strconv"
	"strings"
)

// ZXY describes tile coordinates. Z: Zoom level, X, Y: X and Y coordinates.
type ZXY struct {
	Z, X, Y int
}

func (t ZXY) String() string {
	return fmt.Sprintf("Tile{Z:%v X:%v Y:%v}", t.Z, t.X, t.Y)
}

// ConvertToLangLong converts z, x, y coordinates to latitude and longtitude
/*func (t ZXY) ConvertToLangLong() LatLong {
	var lat, long float64
	n := math.Pi - 2.0*math.Pi*float64(t.Y)/math.Exp2(float64(t.Z))
	lat = 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	long = float64(t.X)/math.Exp2(float64(t.Z))*360.0 - 180.0
	return LatLong{lat, long}
}*/

// ConvertToMeta converts z, x, y to metatiles coordinates.
func (t ZXY) ConvertToMeta() Metatile {
	h := xyToMeta(t.X, t.Y)
	return Metatile{Z: t.Z, Hashes: h}
}

// MinMetaXY returns mininal x and y coordinates contains in metatile.
func (t ZXY) MinMetaXY() (x, y int) {
	mask := MaxMetatileSize - 1
	x = t.X & ^mask
	y = t.Y & ^mask
	return x, y
}

// Path returns filepath of tile, based on Z, X, Y coordinates. Delimiter is "/".
func (t ZXY) Path() string {
	return strconv.Itoa(t.Z) + "/" + strconv.Itoa(t.X) + "/" + strconv.Itoa(t.Y) + ".png"
}

// Path returns filepath of tile
// it is slow implemetation: filepath.Clean is to slow with many dirs elements
/*func (t ZXY) Path(dirs ...string) string {
	p := []string{
		strconv.Itoa(t.Z),
		strconv.Itoa(t.X),
		strconv.Itoa(t.Y) + ".png",
	}
	dirs = append(dirs, p...)
	return filepath.Join(dirs...)
}*/

// ZXYMinURLPathItems is the minimum url items, splitted by separator "/".
const ZXYMinURLPathItems int = 4

// NewZXYFromURL extracts ZXY, style, format from url string.
func NewZXYFromURL(url string) (zxy ZXY, style, format string, err error) {
	items := strings.Split(url, "/")
	il := len(items)
	if il < ZXYMinURLPathItems {
		err = fmt.Errorf("NewZXYFromURL: Wrong url items length: expected %v, got %v", ZXYMinURLPathItems, url)
		return
	}

	// processing -1 value (y.format)
	yformat := strings.Split(items[il-1], ".")
	if len(yformat) != 2 {
		err = fmt.Errorf("NewZXYFromURL: Wrong filename format: %v", items[il-1])
		return
	}

	format = yformat[1]

	zxy.Y, err = strconv.Atoi(yformat[0])
	if err != nil {
		err = fmt.Errorf("NewZXYFromURL: Y: %v", err)
		return
	}

	// processing -2 value (x)
	zxy.X, err = strconv.Atoi(items[il-2])
	if err != nil {
		err = fmt.Errorf("NewZXYFromURL: X: %v", err)
		return
	}

	// processing -3 value (z)
	zxy.Z, err = strconv.Atoi(items[il-3])
	if err != nil {
		err = fmt.Errorf("NewZXYFromURL: Z: %v", err)
		return
	}

	// processing -4 value (style)
	style = items[il-4]

	return
}
