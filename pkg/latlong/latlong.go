package latlong

import "fmt"

// LatLong describes tile coordinates in latitude and longitude format.
type LatLong struct {
	Lat, Long float64
}

func (l LatLong) String() string {
	return fmt.Sprintf("LatLong{%v-%v}", l.Lat, l.Long)
}
