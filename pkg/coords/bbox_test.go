package coords

import (
	"reflect"
	"testing"
)

func TestNewBBoxFromLatLong(t *testing.T) {
	top := LatLong{39.6482, 44.7653}
	bottom := LatLong{25.2919, 61.4949}
	bbox := []Tile{
		Tile{Zoom: 5, X: 19, Y: 12, Ext: "png"},
		Tile{Zoom: 5, X: 19, Y: 13, Ext: "png"},
		Tile{Zoom: 5, X: 20, Y: 12, Ext: "png"},
		Tile{Zoom: 5, X: 20, Y: 13, Ext: "png"},
		Tile{Zoom: 5, X: 21, Y: 12, Ext: "png"},
		Tile{Zoom: 5, X: 21, Y: 13, Ext: "png"},
	}

	result := []Tile{}

	b := NewBBoxFromLatLong([]int{5}, top, bottom, "png")
	for t := range b {
		result = append(result, t)
	}

	if !reflect.DeepEqual(bbox, result) {
		t.Errorf("NewBBoxFromLatLong: expected %q, got %q", bbox, result)
	}
}
