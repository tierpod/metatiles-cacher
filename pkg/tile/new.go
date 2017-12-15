package tile

import (
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/tierpod/metatiles-cacher/pkg/latlong"
)

var reTile = regexp.MustCompile(`(\w+)/(\d+)/(\d+)/(\d+)(\.\w+)`)

// NewFromURL creates Tile from url.
func NewFromURL(url string) (Tile, error) {
	items := reTile.FindStringSubmatch(url)
	if len(items) == 0 {
		return Tile{}, fmt.Errorf("could not parse url string to Tile struct")
	}

	// we can ignore errors because regexp contains `\d+`
	zoom, _ := strconv.Atoi(items[2])
	x, _ := strconv.Atoi(items[3])
	y, _ := strconv.Atoi(items[4])

	return Tile{
		Map:  items[1],
		Zoom: zoom,
		X:    x,
		Y:    y,
		Ext:  items[5],
	}, nil
}

// NewFromLatLong creates Tile from LatLong with zoom and extension.
func NewFromLatLong(l latlong.LatLong, zoom int) Tile {
	var x, y int
	x = int(math.Floor((l.Long + 180.0) / 360.0 * (math.Exp2(float64(zoom)))))
	y = int(math.Floor((1.0 - math.Log(math.Tan(l.Lat*math.Pi/180.0)+1.0/math.Cos(l.Lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(zoom)))))

	return Tile{
		Zoom: zoom,
		X:    x,
		Y:    y,
	}
}
