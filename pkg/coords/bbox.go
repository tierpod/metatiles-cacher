package coords

// NewBBoxFromLatLong create output channel with ZXY tiles
func NewBBoxFromLatLong(zooms []int, top LatLong, bottom LatLong) <-chan (ZXY) {
	bboxChan := make(chan ZXY)

	go func() {
		defer close(bboxChan)
		for _, z := range zooms {
			tTop := top.ConvertToZXY(z)
			tBottom := bottom.ConvertToZXY(z)
			for x := tTop.X; x <= tBottom.X; x++ {
				for y := tTop.Y; y <= tBottom.Y; y++ {
					t := ZXY{Z: z, X: x, Y: y}
					bboxChan <- t
				}
			}
		}
	}()

	return bboxChan
}
