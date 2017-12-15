package metatile

import "fmt"

func ExampleMetatile_Filepath() {
	hashes := [5]int{1, 2, 3, 4, 5}
	mt := Metatile{Zoom: 10, Hashes: hashes, Map: "mapname", X: 0, Y: 0}

	filepath := mt.Filepath("")
	fmt.Println(filepath)

	filepath = mt.Filepath("/var/lib/mod_tile")
	fmt.Println(filepath)

	// Output:
	// mapname/10/5/4/3/2/1.meta
	// /var/lib/mod_tile/mapname/10/5/4/3/2/1.meta
}
