package tile

import (
	"fmt"
	"path"
	"strconv"
)

// Tile describes tile coordinates. zoom level, x and y coordinates, extension, mapname.
type Tile struct {
	Zoom int
	X, Y int
	Ext  string
	Map  string
}

func (t Tile) String() string {
	return fmt.Sprintf("Tile{Zoom:%v X:%v Y:%v Ext:%v Map:%v}", t.Zoom, t.X, t.Y, t.Ext, t.Map)
}

// Filepath returns tile file path, based on basedir and coordinates.
func (t Tile) Filepath(basedir string) string {
	zoom := strconv.Itoa(t.Zoom)
	x := strconv.Itoa(t.X)
	y := strconv.Itoa(t.Y)
	return path.Join(basedir, t.Map, zoom, x, y, t.Ext)
}
