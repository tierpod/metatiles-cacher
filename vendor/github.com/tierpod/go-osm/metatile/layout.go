package metatile

import (
	"errors"
	"fmt"
	"io"
	"math"
)

var (
	// ErrInvalidIndex is the error returned by Decoder when Tile(x, y) not inside metatile file.
	ErrInvalidIndex = errors.New("decoder: invalid index")
	// ErrEmptyData is the error returned by Decoder when Tile has no data (data has zero length).
	ErrEmptyData = errors.New("decoder: empty data")
)

type metaEntry struct {
	Offset int32
	Size   int32
}

// decode tile data for this entry
func (e metaEntry) decode(r io.ReadSeeker) ([]byte, error) {
	if e.Size == 0 {
		return nil, ErrEmptyData
	}

	_, err := r.Seek(int64(e.Offset), 0)
	if err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}

	buf := make([]byte, e.Size)
	n, err := r.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}

	if int32(n) != e.Size {
		return nil, fmt.Errorf("decode: invalid tile size: %v != %v", n, e.Size)
	}

	return buf, nil
}

type metaLayout struct {
	Magic   []byte
	Count   int32
	X, Y, Z int32
	Index   []metaEntry
}

func (m metaLayout) size() int32 {
	return int32(math.Sqrt(float64(m.Count)))
}

func (m metaLayout) tileIndex(x, y int32) (int32, error) {
	i := (x-m.X)*m.size() + (y - m.Y)
	if i >= m.Count {
		return 0, ErrInvalidIndex
	}

	return i, nil
}
