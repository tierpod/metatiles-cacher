package coords

import "testing"

func TestPolygonContains(t *testing.T) {
	// test small polygon
	polygon := Polygon{
		LatLong{Lat: 0, Long: 0},
		LatLong{Lat: 10, Long: 10},
	}

	result := polygon.Contains(LatLong{Lat: 0, Long: 0})
	if result != false {
		t.Errorf("Polygon.Contains: 2 points. Excepted false, got true")
	}

	// test not closed polygon
	polygon = Polygon{
		LatLong{Lat: 0, Long: 0},
		LatLong{Lat: 0, Long: 10},
		LatLong{Lat: 10, Long: 10},
	}

	result = polygon.Contains(LatLong{Lat: 3, Long: 5})
	if result != true {
		t.Errorf("Polygon.Contains: polygon is not closed. Excepted true, got false")
	}

	// test negative coordinates
	polygon = Polygon{
		LatLong{Lat: 10, Long: -10},
		LatLong{Lat: -10, Long: -10},
		LatLong{Lat: -10, Long: 10},
		LatLong{Lat: 10, Long: 10},
		LatLong{Lat: 10, Long: -10},
	}

	result = polygon.Contains(LatLong{Lat: -5, Long: 5})
	if result != true {
		t.Errorf("Polygon.Contains: polygon is not closed. Excepted true, got false")
	}

	// test real coordinates
	polygon = Polygon{
		LatLong{Lat: 55.4903, Long: 65.2110},
		LatLong{Lat: 55.4066, Long: 65.2275},
		LatLong{Lat: 55.4329, Long: 65.3573},
		LatLong{Lat: 55.4969, Long: 65.3878},
		LatLong{Lat: 55.5169, Long: 65.31131},
		LatLong{Lat: 55.4903, Long: 65.2110},
	}

	testData := []struct {
		ll     LatLong
		result bool
	}{
		// inside
		{
			// 12 2790 1285
			LatLong{Lat: 55.4879, Long: 65.2337},
			true,
		},
		{
			// 12 2790 1286
			LatLong{Lat: 55.4416, Long: 65.3003},
			true,
		},
		{
			// 12 2791 1286
			LatLong{Lat: 55.4753, Long: 65.3614},
			true,
		},
		{
			// 12 2791 1285
			LatLong{Lat: 55.4953, Long: 65.3775},
			true,
		},
		{
			// 12 2791 1286
			LatLong{Lat: 55.4732, Long: 65.3181},
			true,
		},
		// outside
		{
			LatLong{Lat: 55.5340, Long: 65.3116},
			false,
		},
		{
			LatLong{Lat: 55.4327, Long: 65.1952},
			false,
		},
		{
			LatLong{Lat: 55.4250, Long: 65.4812},
			false,
		},
	}

	for _, tt := range testData {
		result = polygon.Contains(tt.ll)
		if result != tt.result {
			t.Errorf("Polygon.Contains: point: %v: excepted %v, got %v", tt.ll, tt.result, result)
		}
	}
}
