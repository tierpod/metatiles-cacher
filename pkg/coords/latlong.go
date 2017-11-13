package coords

import (
	"fmt"
	"math"
)

// LatLong describes tile coordinates in latitude and longitude format.
type LatLong struct {
	Lat, Long float64
}

func (l LatLong) String() string {
	return fmt.Sprintf("LatLong{%v-%v}", l.Lat, l.Long)
}

// ToTile converts deg coordinates (LatLong) to num (Tile).
func (l LatLong) ToTile(zoom int) Tile {
	var x, y int
	x = int(math.Floor((l.Long + 180.0) / 360.0 * (math.Exp2(float64(zoom)))))
	y = int(math.Floor((1.0 - math.Log(math.Tan(l.Lat*math.Pi/180.0)+1.0/math.Cos(l.Lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(zoom)))))
	return Tile{Zoom: zoom, X: x, Y: y}
}
