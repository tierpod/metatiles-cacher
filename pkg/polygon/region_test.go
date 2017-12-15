package polygon

import (
	"testing"

	"github.com/tierpod/metatiles-cacher/pkg/latlong"
)

func TestRegionContains(t *testing.T) {
	polygon1 := Polygon{
		latlong.LatLong{Lat: 10.0, Long: 0.0},
		latlong.LatLong{Lat: 0.0, Long: 0.0},
		latlong.LatLong{Lat: 0.0, Long: 5.0},
		latlong.LatLong{Lat: 10.0, Long: 5.0},
		latlong.LatLong{Lat: 10.0, Long: 0.0},
	}

	polygon2 := Polygon{
		latlong.LatLong{Lat: 10.0, Long: 5.0},
		latlong.LatLong{Lat: 0.0, Long: 5.0},
		latlong.LatLong{Lat: 0.0, Long: 10.0},
		latlong.LatLong{Lat: 10.0, Long: 10.0},
		latlong.LatLong{Lat: 10.0, Long: 5.0},
	}

	region := Region{polygon1, polygon2}

	testData := []struct {
		ll     latlong.LatLong
		result bool
	}{
		// inside
		{
			latlong.LatLong{Lat: 3.0, Long: 9.0},
			true,
		},
		{
			latlong.LatLong{Lat: 6.0, Long: 8.0},
			true,
		},
		// outside
		{
			latlong.LatLong{Lat: 15.0, Long: 15.0},
			false,
		},
		{
			latlong.LatLong{Lat: -1.0, Long: -1.0},
			false,
		},
	}

	for _, tt := range testData {
		result := region.Contains(tt.ll)
		if result != tt.result {
			t.Errorf("Polygon.Contains: point: %v: excepted %v, got %v", tt.ll, tt.result, result)
		}
	}
}
