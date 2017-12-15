package latlong

import "math"

// New creates LatLong point from z, x, y coordinates.
func New(z, x, y int) LatLong {
	var lat, long float64
	n := math.Pi - 2.0*math.Pi*float64(y)/math.Exp2(float64(z))
	lat = 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	long = float64(x)/math.Exp2(float64(z))*360.0 - 180.0
	return LatLong{Lat: lat, Long: long}
}
