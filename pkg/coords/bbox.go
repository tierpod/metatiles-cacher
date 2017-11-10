package coords

// NewBBoxFromLatLong calculates bbox from top to bottom coordinates for each zoom level in zooms.
// Returns output chan with tiles from this bbox.
func NewBBoxFromLatLong(zooms []int, top LatLong, bottom LatLong, ext string) <-chan (Tile) {
	bboxChan := make(chan Tile)

	go func() {
		defer close(bboxChan)
		for _, z := range zooms {
			tTop := top.ToTile(z)
			tBottom := bottom.ToTile(z)
			for x := tTop.X; x <= tBottom.X; x++ {
				for y := tTop.Y; y <= tBottom.Y; y++ {
					t := Tile{Zoom: z, X: x, Y: y, Ext: ext}
					bboxChan <- t
				}
			}
		}
	}()

	return bboxChan
}
