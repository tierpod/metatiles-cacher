package polygon

import "github.com/tierpod/metatiles-cacher/pkg/latlong"

// Region is the slice of closed polygons.
type Region []Polygon

// Contains checks if LatLong point contains in each polygon.
func (r Region) Contains(pt latlong.LatLong) bool {
	for _, p := range r {
		inside := p.Contains(pt)
		if inside {
			return true
		}
	}

	return false
}
