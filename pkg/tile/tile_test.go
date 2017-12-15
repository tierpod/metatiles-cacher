package tile

import "fmt"

func ExampleTile_Filepath() {
	t := Tile{Zoom: 1, X: 1, Y: 1, Ext: ".png", Map: "mapname"}

	tilepath := t.Filepath("")
	fmt.Println(tilepath)

	tilepath = t.Filepath("/var/cache/tile")
	fmt.Println(tilepath)

	// Output:
	// mapname/1/1/1/.png
	// /var/cache/tile/mapname/1/1/1/.png
}
