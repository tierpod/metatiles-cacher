// Package point contains geo points description and conversation.
package point

import (
	"fmt"
	"math"
)

// ZXY represents point with (z, x, y) coordinates.
type ZXY struct {
	Z, X, Y int
}

func (p ZXY) String() string {
	return fmt.Sprintf("{Z:%v X:%v Y:%v}", p.Z, p.X, p.Y)
}

// ToLatLong converts ZXY point to latitude and longtitude coordinates.
func (p ZXY) ToLatLong() LatLong {
	var lat, long float64
	n := math.Pi - 2.0*math.Pi*float64(p.Y)/math.Exp2(float64(p.Z))
	lat = 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	long = float64(p.X)/math.Exp2(float64(p.Z))*360.0 - 180.0
	return LatLong{Lat: lat, Long: long}
}

// LatLong represents point with (Lat, Long) coordinates
type LatLong struct {
	Lat, Long float64
}

func (p LatLong) String() string {
	return fmt.Sprintf("{Lat:%v Long:%v}", p.Lat, p.Long)
}

// ToZXY converts latitude and longtitude for given zoom to (z, x, y) coordinates.
func (p LatLong) ToZXY(zoom int) ZXY {
	var x, y int
	x = int(math.Floor((p.Long + 180.0) / 360.0 * (math.Exp2(float64(zoom)))))
	y = int(math.Floor((1.0 - math.Log(math.Tan(p.Lat*math.Pi/180.0)+1.0/math.Cos(p.Lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(zoom)))))
	return ZXY{Z: zoom, X: x, Y: y}
}
