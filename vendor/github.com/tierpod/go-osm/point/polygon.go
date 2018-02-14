package point

// Polygon is the closed polygon area inside LatLong points as vertices.
type Polygon []LatLong

// Contains checks if LatLong point contains in polygon. Use ray-casting algorithm.
// http://rosettacode.org/wiki/Ray-casting_algorithm#Go
func (p Polygon) Contains(pt LatLong) bool {
	inside := false
	polyLen := len(p)

	if polyLen < 3 {
		return false
	}

	if !(p[0] == p[polyLen-1]) {
		// polygon is not closed, use last point as first point.
		inside = rayIntersectsSegment(pt, p[polyLen-1], p[0])
	}

	for i := 1; i < polyLen; i++ {
		if rayIntersectsSegment(pt, p[i-1], p[i]) {
			inside = !inside
		}
	}
	return inside
}

// lat=x, long=y
func rayIntersectsSegment(p, a, b LatLong) bool {
	return (a.Long > p.Long) != (b.Long > p.Long) &&
		p.Lat < (b.Lat-a.Lat)*(p.Long-a.Long)/(b.Long-a.Long)+a.Lat
}
