package tile

import (
	"fmt"
	"regexp"
	"strconv"
)

var reTile = regexp.MustCompile(`(\w+)/(\d+)/(\d+)/(\d+)(\.\w+)`)

// NewFromURL extracts Tile coordinates, style, format from url string.
func NewFromURL(url string) (t Tile, err error) {
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
