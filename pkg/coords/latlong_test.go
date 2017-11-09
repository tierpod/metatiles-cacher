package coords

import "testing"

func TestLatLongString(t *testing.T) {
	l := LatLong{1.2345, 5.4321}
	result := "LatLong{1.2345-5.4321}"
	if l.String() != result {
		t.Errorf("LatLong String(): expected %v, got %v", result, l.String())
	}
}
