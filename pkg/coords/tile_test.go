package coords

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTileString(t *testing.T) {
	tile := Tile{10, 697, 321}
	result := "Tile{Zoom:10 X:697 Y:321}"

	if tile.String() != result {
		t.Errorf("Tile.String(): expected %v, got %v", result, tile.String())
	}
}

/*func TestConvertToLatLong(t *testing.T) {
	var testData = []struct {
		zxy    Tile
		result LatLong
	}{
		{
			Tile{5, 19, 12},
			LatLong{39.6482, 44.7653},
		},
	}

	for _, tt := range testData {
		ll := tt.zxy.ConvertToLangLong()
		if !reflect.DeepEqual(ll, tt.result) {
			t.Errorf("Tile ConvertToLatLong: expected %v, got %v", tt.result, ll)
		}
	}
}*/

func TestConvertToMeta(t *testing.T) {
	var testData = []struct {
		tile   Tile
		result Metatile
	}{
		{
			Tile{1, 1, 1},
			Metatile{1, [5]int{0, 0, 0, 0, 0}},
		},
		{
			Tile{10, 697, 321},
			Metatile{10, [5]int{128, 180, 33, 0, 0}},
		},
		{
			Tile{17, 89325, 41158},
			Metatile{17, [5]int{128, 236, 192, 90, 16}},
		},
	}

	for _, tt := range testData {
		m := tt.tile.ToMetatile()
		if m != tt.result {
			t.Errorf("Tile ConvertToLatLong: expected %v, got %v", tt.result, m)
		}
	}
}

func BenchmarkTilePath(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t := Tile{0, 0, 0}
		t.Path()
	}
}

func TestNewTileFromURL(t *testing.T) {
	testData := []struct {
		url           string
		style, format string
		tile          Tile
	}{
		{
			"/maps/style/10/697/321.png",
			"style", "png",
			Tile{10, 697, 321},
		},
		{
			"/api/add/style/10/697/321.png",
			"style", "png",
			Tile{10, 697, 321},
		},
		{
			"http://localhost:8080/test/api/add/style/10/697/321.png",
			"style", "png",
			Tile{10, 697, 321},
		},
	}

	for _, tt := range testData {
		tile, style, format, err := NewTileFromURL(tt.url)
		if err != nil {
			t.Errorf("Got error: %v", err)
		}

		if format != tt.format {
			t.Errorf("Expected format %v, got %v", tt.format, format)
		}

		if style != tt.style {
			t.Errorf("Expected format %v, got %v", tt.style, style)
		}

		if !reflect.DeepEqual(tile, tt.tile) {
			t.Errorf("Expected %v, got %v", tt.tile, tile)
		}
	}

	testDataErr := []struct {
		url    string
		errStr string
	}{
		{
			"/maps",
			"NewTileFromURL: wrong url items length: /maps (2/5)",
		},
		{
			"/maps/style",
			"NewTileFromURL: wrong url items length: /maps/style (3/5)",
		},
		{
			"/maps/style/z/1/1.png",
			"NewTileFromURL: Z: strconv.Atoi: parsing \"z\": invalid syntax",
		},
		{
			"/maps/style/1/x/1.png",
			"NewTileFromURL: X: strconv.Atoi: parsing \"x\": invalid syntax",
		},
		{
			"/maps/style/1/1/y.png",
			"NewTileFromURL: Y: strconv.Atoi: parsing \"y\": invalid syntax",
		},
		{
			"/maps/style/1/1/y_png",
			"NewTileFromURL: wrong filename format: y_png",
		},
	}

	for _, tt := range testDataErr {
		_, _, _, err := NewTileFromURL(tt.url)
		if err.Error() != tt.errStr {
			t.Errorf("Expected error: %v, got: %v", tt.errStr, err)
		}
	}
}

func TestTileMinMetaXY(t *testing.T) {
	testData := []struct {
		tile Tile
		x, y int
	}{
		{
			Tile{1, 1, 1},
			0, 0,
		},
		{
			Tile{10, 697, 321},
			696, 320,
		},
		{
			Tile{17, 89325, 41158},
			89320, 41152,
		},
	}

	for _, tt := range testData {
		x, y := tt.tile.MinMetatileXY()
		if x != tt.x || y != tt.y {
			t.Errorf("MinMetaXY: expected {X:%v Y:%v}, got {X:%v Y:%v}", tt.x, tt.y, x, y)
		}
	}
}

func BenchmarkNewTileFromURL(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewTileFromURL("/maps/style/1/2/3.png")
	}
}

func ExampleNewTileFromURL() {
	t, style, format, _ := NewTileFromURL("/maps/style/10/2/3.png")

	fmt.Printf("%v %v %v", t, style, format)

	// Output:
	// Tile{Zoom:10 X:2 Y:3} style png
}

func ExampleTile_Path() {
	t := Tile{Zoom: 10, X: 2, Y: 3}

	fmt.Printf("%v\n", t.Path())

	// Output:
	// 10/2/3.png
}
