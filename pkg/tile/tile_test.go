package tile

import "fmt"

func ExampleNewFromURL() {
	urls := []string{
		"map/1/2/3.png",
		"http://localhost:8080/map/1/2/3.png",
		"1/2/3.png",
		"map/z/x/y.png",
	}

	for _, url := range urls {
		t, err := NewFromURL(url)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		fmt.Println(t)
	}

	// Output:
	// Tile{Zoom:1 X:2 Y:3 Ext:.png Map:map}
	// Tile{Zoom:1 X:2 Y:3 Ext:.png Map:map}
	// error: could not parse url string to Tile struct
	// error: could not parse url string to Tile struct
}

func ExampleTile_Filepath() {
	t := New(10, 1, 2, ".png", "map")
	fmt.Println(t.Filepath("/var/cache/tiles"))

	// Output:
	// /var/cache/tiles/map/10/1/2.png
}
