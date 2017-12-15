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

func ExampleMetatile_Size() {
	zooms := []int{1, 2, 3, 8}
	for _, zoom := range zooms {
		mt := Metatile{Zoom: zoom}
		fmt.Println(zoom, mt.Size())
	}

	// Output:
	// 1 2
	// 2 4
	// 3 8
	// 8 8
}

func ExampleMetatile_XYBox() {
	url := "/var/lib/mod_tile/map/10/0/0/33/180/128.meta"
	mt, _ := NewFromURL(url)
	fmt.Println("X:", mt.XYBox().X)
	fmt.Println("Y:", mt.XYBox().Y)

	// Output:
	// X: [696 697 698 699 700 701 702 703]
	// Y: [320 321 322 323 324 325 326 327]
}

func ExampleXYOffset() {
	xx := []int{0, 1}
	yy := []int{0, 1}

	for x := range xx {
		for y := range yy {
			offset := XYOffset(x, y)
			fmt.Printf("(%v, %v): %v\n", x, y, offset)
		}
	}

	// Output:
	// (0, 0): 0
	// (0, 1): 1
	// (1, 0): 8
	// (1, 1): 9
}
