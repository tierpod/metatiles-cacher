package point

// Region includes many polygons.
type Region []Polygon

// Contains checks if LatLong point contains in one of the polygon inside region.
func (r Region) Contains(pt LatLong) bool {
	for _, p := range r {
		inside := p.Contains(pt)
		if inside {
			return true
		}
	}

	return false
}
