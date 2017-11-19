package coords

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Tile describes tile coordinates. Zoom level, X, Y and extension.
type Tile struct {
	Zoom int
	X, Y int
	Ext  string
}

func (t Tile) String() string {
	return fmt.Sprintf("Tile{Zoom:%v X:%v Y:%v Ext:%v}", t.Zoom, t.X, t.Y, t.Ext)
}

// ToLangLong converts z, x, y coordinates to latitude and longitude.
func (t Tile) ToLangLong() LatLong {
	var lat, long float64
	n := math.Pi - 2.0*math.Pi*float64(t.Y)/math.Exp2(float64(t.Zoom))
	lat = 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	long = float64(t.X)/math.Exp2(float64(t.Zoom))*360.0 - 180.0
	return LatLong{Lat: lat, Long: long}
}

// ToMetatile converts Tile to Metatiles coordinates.
func (t Tile) ToMetatile() Metatile {
	h := xyToMetatile(t.X, t.Y)
	return Metatile{Zoom: t.Zoom, Hashes: h}
}

func xyToMetatile(x, y int) hashes {
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

// MinMetatileXY calculates min x and y coordinates that can be stored in metatile with MaxMetatileSize size.
func (t Tile) MinMetatileXY() (x, y int) {
	mask := MaxMetatileSize - 1
	x = t.X & ^mask
	y = t.Y & ^mask
	return x, y
}

// Path returns filepath of tile, based on Z, X, Y coordinates. Delimiter is "/".
func (t Tile) Path() string {
	return strconv.Itoa(t.Zoom) + "/" + strconv.Itoa(t.X) + "/" + strconv.Itoa(t.Y) + `.` + t.Ext
}

// Mimetype return mimetype based on tile extension.
func (t Tile) Mimetype() (string, error) {
	switch t.Ext {
	case "png":
		return "image/png", nil
	case "json", "topojson", "geojson":
		return "application/json", nil
	case "mvt":
		return "application/vnd.mapbox-vector-tile", nil
	default:
		return "", fmt.Errorf("unknown mimetype for extension \"%v\"", t.Ext)
	}
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

// TileMinURLPathItems is the minimum url items length, splitted by separator "/".
//
// Example: /style/zoom/x/y/y.png has length 5.
const TileMinURLPathItems int = 5

// NewTileFromURL extracts Tile coordinates, style, format from url string.
func NewTileFromURL(url string) (t Tile, style string, err error) {
	items := strings.Split(url, "/")
	il := len(items)
	if il < TileMinURLPathItems {
		err = fmt.Errorf("NewTileFromURL: wrong url items length: %v (%v/%v)", url, il, TileMinURLPathItems)
		return
	}

	// processing -1 value (y.png), split to ["y" "png"]
	yext := strings.Split(items[il-1], ".")
	if len(yext) != 2 {
		err = fmt.Errorf("NewTileFromURL: wrong filename format: %v", items[il-1])
		return
	}

	t.Ext = yext[1]

	t.Y, err = strconv.Atoi(yext[0])
	if err != nil {
		err = fmt.Errorf("NewTileFromURL: Y: %v", err)
		return
	}

	// processing -2 value (x)
	t.X, err = strconv.Atoi(items[il-2])
	if err != nil {
		err = fmt.Errorf("NewTileFromURL: X: %v", err)
		return
	}

	// processing -3 value (z)
	t.Zoom, err = strconv.Atoi(items[il-3])
	if err != nil {
		err = fmt.Errorf("NewTileFromURL: Z: %v", err)
		return
	}

	// processing -4 value (style)
	style = items[il-4]

	return
}
