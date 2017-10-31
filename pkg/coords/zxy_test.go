package coords

import (
	"fmt"
	"reflect"
	"testing"
)

func TestZXYString(t *testing.T) {
	zxy := ZXY{10, 697, 321}
	result := "Tile{Z:10 X:697 Y:321}"

	if zxy.String() != result {
		t.Errorf("ZXY String(): expected %v, got %v", result, zxy.String())
	}
}

/*func TestConvertToLatLong(t *testing.T) {
	var testData = []struct {
		zxy    ZXY
		result LatLong
	}{
		{
			ZXY{5, 19, 12},
			LatLong{39.6482, 44.7653},
		},
	}

	for _, tt := range testData {
		ll := tt.zxy.ConvertToLangLong()
		if !reflect.DeepEqual(ll, tt.result) {
			t.Errorf("ZXY ConvertToLatLong: expected %v, got %v", tt.result, ll)
		}
	}
}*/

func TestConvertToMeta(t *testing.T) {
	var testData = []struct {
		zxy    ZXY
		result Metatile
	}{
		{
			ZXY{1, 1, 1},
			Metatile{1, [5]int{0, 0, 0, 0, 0}},
		},
		{
			ZXY{10, 697, 321},
			Metatile{10, [5]int{128, 180, 33, 0, 0}},
		},
		{
			ZXY{17, 89325, 41158},
			Metatile{17, [5]int{128, 236, 192, 90, 16}},
		},
	}

	for _, tt := range testData {
		m := tt.zxy.ConvertToMeta()
		if m != tt.result {
			t.Errorf("ZXY ConvertToLatLong: expected %v, got %v", tt.result, m)
		}
	}
}

func BenchmarkZXYPath1000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t := ZXY{0, 0, 0}
		t.Path()
	}
}

func TestNewZXYFromURL(t *testing.T) {
	testData := []struct {
		url           string
		style, format string
		zxy           ZXY
	}{
		{
			"/maps/style/10/697/321.png",
			"style", "png",
			ZXY{10, 697, 321},
		},
		{
			"/api/add/style/10/697/321.png",
			"style", "png",
			ZXY{10, 697, 321},
		},
		{
			"http://localhost:8080/test/api/add/style/10/697/321.png",
			"style", "png",
			ZXY{10, 697, 321},
		},
	}

	for _, tt := range testData {
		zxy, style, format, err := NewZXYFromURL(tt.url)
		if err != nil {
			t.Errorf("Got error: %v", err)
		}

		if format != tt.format {
			t.Errorf("Expected format %v, got %v", tt.format, format)
		}

		if style != tt.style {
			t.Errorf("Expected format %v, got %v", tt.style, style)
		}

		if !reflect.DeepEqual(zxy, tt.zxy) {
			t.Errorf("Expected %v, got %v", tt.zxy, zxy)
		}
	}

	testDataErr := []struct {
		url    string
		errStr string
	}{
		{
			"/maps/",
			"NewZXYFromURL: Wrong url items length: expected 4, got /maps/",
		},
		{
			"/maps/style",
			"NewZXYFromURL: Wrong url items length: expected 4, got /maps/style",
		},
		{
			"/maps/style/z/1/1.png",
			"NewZXYFromURL: Z: strconv.Atoi: parsing \"z\": invalid syntax",
		},
		{
			"/maps/style/1/x/1.png",
			"NewZXYFromURL: X: strconv.Atoi: parsing \"x\": invalid syntax",
		},
		{
			"/maps/style/1/1/y.png",
			"NewZXYFromURL: Y: strconv.Atoi: parsing \"y\": invalid syntax",
		},
		{
			"/maps/style/1/1/y_png",
			"NewZXYFromURL: Wrong filename format: y_png",
		},
	}

	for _, tt := range testDataErr {
		_, _, _, err := NewZXYFromURL(tt.url)
		if err.Error() != tt.errStr {
			t.Errorf("Expected error: %v, got: %v", tt.errStr, err)
		}
	}
}

func TestZXYMinMetaXY(t *testing.T) {
	testData := []struct {
		zxy  ZXY
		x, y int
	}{
		{
			ZXY{1, 1, 1},
			0, 0,
		},
		{
			ZXY{10, 697, 321},
			696, 320,
		},
		{
			ZXY{17, 89325, 41158},
			89320, 41152,
		},
	}

	for _, tt := range testData {
		x, y := tt.zxy.MinMetaXY()
		if x != tt.x || y != tt.y {
			t.Errorf("MinMetaXY: expected {X:%v Y:%v}, got {X:%v Y:%v}", tt.x, tt.y, x, y)
		}
	}
}

func ExampleNewZXYFromURL() {
	zxy, style, format, _ := NewZXYFromURL("/maps/style/1/2/3.png")

	fmt.Printf("%v %v %v", zxy, style, format)

	// Output:
	// Tile{Z:1 X:2 Y:3} style png
}

func ExampleZXY_Path() {
	zxy := ZXY{Z: 1, X: 2, Y: 3}

	fmt.Printf("%v\n", zxy.Path())

	// Output:
	// 1/2/3.png
}
