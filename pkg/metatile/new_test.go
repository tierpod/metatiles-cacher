package metatile

import (
	"fmt"

	"github.com/tierpod/metatiles-cacher/pkg/tile"
)

func ExampleNewFromURL() {
	urls := []string{
		"map/10/0/0/33/180/128.meta",
		"/var/lib/mod_tile/map/10/0/0/33/180/128.meta",
		"http://localhost:8080/maps/map/10/0/0/33/180/128.meta",
		"map/ZOOM/0/0/33/180/128.meta",
	}

	for _, url := range urls {
		mt, err := NewFromURL(url)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		fmt.Println(mt)
	}

	// Output:
	// Metatile{Zoom:10 Hashes:[128 180 33 0 0] Map:map Ext:.meta X:696 Y:320}
	// Metatile{Zoom:10 Hashes:[128 180 33 0 0] Map:map Ext:.meta X:696 Y:320}
	// Metatile{Zoom:10 Hashes:[128 180 33 0 0] Map:map Ext:.meta X:696 Y:320}
	// error: could not parse url string to Metatile struct
}

func ExampleNewFromTile() {
	t := tile.Tile{Zoom: 10, X: 697, Y: 321, Ext: ".png"}
	mt := NewFromTile(t)
	fmt.Println(mt)
	fmt.Println(mt.Filepath("/var/lib/mod_tile"))

	// Output:
	// Metatile{Zoom:10 Hashes:[128 180 33 0 0] Map: Ext:.meta X:696 Y:320}
	// /var/lib/mod_tile/10/0/0/33/180/128.meta
}
