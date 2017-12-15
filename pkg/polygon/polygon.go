package polygon

import "github.com/tierpod/metatiles-cacher/pkg/latlong"

// Polygon is the closed polygon with LatLong points as vertices.
type Polygon []latlong.LatLong

// Contains checks if LatLong point contains in polygon. Use ray-casting algorithm.
// http://rosettacode.org/wiki/Ray-casting_algorithm#Go
func (p Polygon) Contains(pt latlong.LatLong) bool {
	in := false
	pl := len(p)

	if pl < 3 {
		return false
	}

	if !(p[0] == p[pl-1]) {
		// polygon is not closed, use last point as first point.
		in = rayIntersectsSegment(pt, p[pl-1], p[0])
	}

	for i := 1; i < pl; i++ {
		if rayIntersectsSegment(pt, p[i-1], p[i]) {
			in = !in
		}
	}
	return in
}

// lat=x, long=y
func rayIntersectsSegment(p, a, b latlong.LatLong) bool {
	return (a.Long > p.Long) != (b.Long > p.Long) &&
		p.Lat < (b.Lat-a.Lat)*(p.Long-a.Long)/(b.Long-a.Long)+a.Lat
}
