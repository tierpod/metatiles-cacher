package util

import (
	"fmt"
	"testing"
)

func TestMakeIntSlice(t *testing.T) {
	result := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := MakeIntSlice(1, 10)

	if len(r) != len(result) {
		t.Errorf("MakeIntSlice(1, 17): invalid result slice length (expected %v, got %v)", len(result), len(r))
	}

	for i := range r {
		if r[i] != result[i] {
			t.Errorf("MakeIntSlice(1, 17): invalid slice item (expected: %v, got %v)", result[i], r[i])
		}
	}
}

func ExampleMakeIntSlice() {
	r := MakeIntSlice(1, 10)
	fmt.Printf("%v\n", r)

	// Output:
	// [1 2 3 4 5 6 7 8 9]
}

func ExampleDigestString() {
	d := DigestString("teststring")
	fmt.Printf("%v\n", d)

	// Output:
	// d67c5cbf5b01c9f91932e3b8def5e5f8
}

func ExampleMimetype() {
	exts := []string{".png", ".json", ".topojson", ".geojson", ".mvt", ".unknown"}

	for _, ext := range exts {
		mt, err := Mimetype(ext)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		fmt.Println(ext, mt)
	}

	// Output:
	// .png image/png
	// .json application/json
	// .topojson application/json
	// .geojson application/json
	// .mvt application/vnd.mapbox-vector-tile
	// error: unknown mimetype for extension ".unknown"
}
