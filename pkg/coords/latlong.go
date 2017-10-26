package coords

import (
	"fmt"
	"math"
)

// LatLong is basic structure with latitude and longtitude coordinates
type LatLong struct {
	Lat, Long float64
}

func (t LatLong) String() string {
	return fmt.Sprintf("LatLong{%v-%v}", t.Lat, t.Long)
}

// ConvertToZXY converts deg to num
func (t LatLong) ConvertToZXY(zoom int) ZXY {
	var x, y int
	x = int(math.Floor((t.Long + 180.0) / 360.0 * (math.Exp2(float64(zoom)))))
	y = int(math.Floor((1.0 - math.Log(math.Tan(t.Lat*math.Pi/180.0)+1.0/math.Cos(t.Lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(zoom)))))
	return ZXY{zoom, x, y}
}
