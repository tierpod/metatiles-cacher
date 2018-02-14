package metatile

type hashes [5]int

func (h hashes) XY() (x, y int) {
	for i := 4; i >= 0; i-- {
		x <<= 4
		y <<= 4
		x = x | (h[i]&0xf0)>>4
		y = y | (h[i] & 0x0f)
	}

	return x, y
}

func xyToHashes(x, y int) hashes {
	var xx, yy, mask int

	mask = MaxSize - 1
	xx = x & ^mask
	yy = y & ^mask
	h := hashes{}

	for i := 0; i < 5; i++ {
		h[i] = ((xx & 0x0f) << 4) | (yy & 0x0f)
		xx >>= 4
		yy >>= 4
	}

	return h
}
