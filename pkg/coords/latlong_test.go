package coords

import "testing"

func TestLatLongString(t *testing.T) {
	ll := LatLong{1.2345, 5.4321}
	result := "LatLong{1.2345-5.4321}"
	if ll.String() != result {
		t.Errorf("LatLong String(): expected %v, got %v", result, ll.String())
	}
}

/*func TestConvertToZXY(t *testing.T) {
}*/
