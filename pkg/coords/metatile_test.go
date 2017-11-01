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
	testData := []struct {
		m      Metatile
		result XYBox
	}{
		{
			Metatile{10, [5]int{128, 180, 33, 0, 0}},
			XYBox{
				[]int{696, 697, 698, 699, 700, 701, 702, 703},
				[]int{320, 321, 322, 323, 324, 325, 326, 327},
			},
		},
		{
			Metatile{1, [5]int{0, 0, 0, 0, 0}},
			XYBox{
				[]int{0, 1},
				[]int{0, 1},
			},
		},
	}

	for _, tt := range testData {
		xybox := tt.m.ConvertToXYBox()
		if !reflect.DeepEqual(xybox, tt.result) {
			t.Errorf("ConvertToXYBox: expected %v, got %v", tt.result, xybox)
		}
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

func TestNewMetaFromURL(t *testing.T) {
	testData := []struct {
		url   string
		style string
		meta  Metatile
	}{
		{
			"/metatile/path/style/10/128/0/0/0/0.meta",
			"style",
			Metatile{Z: 10, Hashes: [5]int{0, 0, 0, 0, 128}},
		},
		{
			"/style/10/128/0/0/0/33.meta",
			"style",
			Metatile{Z: 10, Hashes: [5]int{33, 0, 0, 0, 128}},
		},
		{
			"http://localhost/test/add/style/12/128/0/128/0/33.meta",
			"style",
			Metatile{Z: 12, Hashes: [5]int{33, 0, 128, 0, 128}},
		},
	}

	for _, tt := range testData {
		meta, style, err := NewMetaFromURL(tt.url)
		if err != nil {
			t.Errorf("Got error: %v", err)
		}

		if style != tt.style {
			t.Errorf("Expected format %v, got %v", tt.style, style)
		}

		if !reflect.DeepEqual(meta, tt.meta) {
			t.Errorf("Expected %v, got %v", tt.meta, meta)
		}
	}

	testDataErr := []struct {
		url    string
		errStr string
	}{
		{
			"/maps/",
			"NewMetaFromURL: Wrong url items length: expected 7, got /maps/",
		},
		{
			"/maps/style/z/4/3/2/1/0.meta",
			"NewMetaFromURL: Z: strconv.Atoi: parsing \"z\": invalid syntax",
		},
		{
			"/maps/style/10/h4/3/2/1/0.meta",
			"NewMetaFromURL: h4: strconv.Atoi: parsing \"h4\": invalid syntax",
		},
		{
			"/maps/style/10/4/h3/2/1/0.meta",
			"NewMetaFromURL: h3: strconv.Atoi: parsing \"h3\": invalid syntax",
		},
		{
			"/maps/style/10/4/3/h2/1/0.meta",
			"NewMetaFromURL: h2: strconv.Atoi: parsing \"h2\": invalid syntax",
		},
		{
			"/maps/style/10/4/3/2/h1/0.meta",
			"NewMetaFromURL: h1: strconv.Atoi: parsing \"h1\": invalid syntax",
		},
		{
			"/maps/style/10/4/3/2/1/h0.meta",
			"NewMetaFromURL: h0: strconv.Atoi: parsing \"h0\": invalid syntax",
		},
		{
			"/maps/style/10/4/3/2/1/h0",
			"NewMetaFromURL: Wrong last item: h0",
		},
	}

	for _, tt := range testDataErr {
		_, _, err := NewMetaFromURL(tt.url)
		if err.Error() != tt.errStr {
			t.Errorf("Expected error: %v, got: %v", tt.errStr, err)
		}
	}
}

func BenchmarkNewMetaFromURL(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewMetaFromURL("/maps/style/10/0/1/2/3/4.meta")
	}
}

func ExampleMetatile_Path() {
	meta := Metatile{Z: 1, Hashes: [5]int{2, 3, 4, 5, 6}}

	fmt.Printf("%v\n", meta.Path())

	// Output:
	// 1/6/5/4/3/2.meta
}

func ExampleNewMetaFromURL() {
	meta, style, _ := NewMetaFromURL("/metatile/path/style/10/128/0/0/0/0.meta")
	fmt.Printf("%v %v\n", meta, style)

	// Output:
	// Metatile{Z:10 Hashes:[0 0 0 0 128]} style
}
