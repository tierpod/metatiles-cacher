package coords

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMetatileString(t *testing.T) {
	m := Metatile{10, [5]int{128, 180, 33, 0, 0}}
	result := "Metatile{Z:10 Hashes:[128 180 33 0 0]}"

	if m.String() != result {
		t.Errorf("Metatile String(): expected %v, got %v", result, m.String())
	}
}

func TestConvertToXYBox(t *testing.T) {
	m := Metatile{10, [5]int{128, 180, 33, 0, 0}}
	xybox := m.ConvertToXYBox()
	result := XYBox{
		[]int{696, 697, 698, 699, 700, 701, 702, 703},
		[]int{320, 321, 322, 323, 324, 325, 326, 327},
	}

	if !reflect.DeepEqual(xybox, result) {
		t.Errorf("ConvertToXYBox: expected %v, got %v", result, xybox)
	}
}

func TestMetatileSize(t *testing.T) {
	testData := []struct {
		size   int
		result int
	}{
		{1, 2},
		{2, 4},
		{3, 8},
		{8, 8},
	}

	for _, tt := range testData {
		s := metatileSize(tt.size)
		if s != tt.result {
			t.Errorf("ZXY metatileSize: expected %v, got %v", tt.result, s)
		}
	}
}

func BenchmarkMetatilePath1000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		m := Metatile{Z: 1, Hashes: [5]int{0, 0, 0, 0, 128}}
		m.Path()
	}
}

func TestMetatileMinXY(t *testing.T) {
	testData := []struct {
		meta Metatile
		x, y int
	}{
		{
			Metatile{Z: 1, Hashes: [5]int{0, 0, 0, 0, 0}},
			0, 0,
		},
		{
			Metatile{Z: 10, Hashes: [5]int{128, 180, 33, 0, 0}},
			696, 320,
		},
		{
			Metatile{Z: 17, Hashes: [5]int{128, 236, 192, 90, 16}},
			89320, 41152,
		},
	}

	for _, tt := range testData {
		x, y := tt.meta.MinXY()
		if x != tt.x || y != tt.y {
			t.Errorf("MinMetaXY: expected {X:%v Y:%v}, got {X:%v Y:%v}", tt.x, tt.y, x, y)
		}
	}
}

func ExampleMetatile_Path() {
	meta := Metatile{Z: 1, Hashes: [5]int{2, 3, 4, 5, 6}}

	fmt.Printf("%v\n", meta.Path())

	// Output:
	// 1/6/5/4/3/2.meta
}
