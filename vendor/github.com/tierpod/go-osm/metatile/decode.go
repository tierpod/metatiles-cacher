package metatile

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/tierpod/go-osm/point"
)

// Decoder is the metatile file decoder wrapper.
type Decoder struct {
	ml *metaLayout
	r  io.ReadSeeker
}

// NewDecoder reads metatile from r, parses layout and returns Decoder.
func NewDecoder(r io.ReadSeeker) (*Decoder, error) {
	ml, err := decodeHeader(r)
	if err != nil {
		return nil, err
	}

	return &Decoder{
		ml: ml,
		r:  r,
	}, nil
}

// Header returns metatile header: x, y, z coordinates and count of tiles.
func (m *Decoder) Header() (x, y, z, count int32) {
	return m.ml.X, m.ml.Y, m.ml.Z, m.ml.Count
}

// Size returns metatile size.
func (m *Decoder) Size() int {
	return int(m.ml.size())
}

// Entries is the array of metaEntry.
type Entries []metaEntry

// Entries returns metatile index table (offsets and sizes).
func (m *Decoder) Entries() Entries {
	return m.ml.Index
}

// Tile reads data for tile with (x, y) coordinates.
func (m *Decoder) Tile(x, y int) ([]byte, error) {
	i, err := m.ml.tileIndex(int32(x), int32(y))
	if err != nil {
		return nil, err
	}

	entry := m.ml.Index[i]
	data, err := entry.decode(m.r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Tiles reads data for all tiles in metatile and returns all data as array of data (includes empty
// data).
func (m *Decoder) Tiles() ([][]byte, error) {
	var tiles [][]byte

	for _, entry := range m.ml.Index {
		data, err := entry.decode(m.r)
		if err != nil && err != ErrEmptyData {
			return nil, err
		}

		tiles = append(tiles, data)
	}

	return tiles, nil
}

// TilesMap reads data for all tiles in metatile and returns only none-empty data as map.
func (m *Decoder) TilesMap() (map[point.ZXY][]byte, error) {
	r := make(map[point.ZXY][]byte)

	for i, entry := range m.ml.Index {
		data, err := entry.decode(m.r)
		if err != nil && err != ErrEmptyData {
			return nil, err
		}

		if len(data) == 0 {
			continue
		}

		x, y := IndexToXY(i)
		p := point.ZXY{Z: int(m.ml.Z), X: int(m.ml.X) + x, Y: int(m.ml.Y) + y}
		r[p] = data
	}

	return r, nil
}

// decodeHeader reads metatile from r and decodes header to metaLayout struct.
func decodeHeader(r io.Reader) (*metaLayout, error) {
	endian := binary.LittleEndian
	ml := new(metaLayout)

	ml.Magic = make([]byte, 4)
	err := binary.Read(r, endian, &ml.Magic)
	if err != nil {
		return nil, err
	}
	if ml.Magic[0] != 'M' || ml.Magic[1] != 'E' || ml.Magic[2] != 'T' || ml.Magic[3] != 'A' {
		return nil, fmt.Errorf("invalid Magic field: %v", ml.Magic)
	}

	if err = binary.Read(r, endian, &ml.Count); err != nil {
		return nil, err
	}
	if err = binary.Read(r, endian, &ml.X); err != nil {
		return nil, err
	}
	if err = binary.Read(r, endian, &ml.Y); err != nil {
		return nil, err
	}
	if err = binary.Read(r, endian, &ml.Z); err != nil {
		return nil, err
	}

	for i := int32(0); i < ml.Count; i++ {
		var entry metaEntry
		if err = binary.Read(r, endian, &entry); err != nil {
			return nil, err
		}
		ml.Index = append(ml.Index, entry)
	}

	return ml, nil
}

// XYToIndex returns offset of tile data inside metatile.
func XYToIndex(x, y int) int {
	mask := MaxSize - 1
	return (x&mask)*MaxSize + (y & mask)
}

// IndexToXY returns (x, y) coordinates from inex.
func IndexToXY(i int) (x, y int) {
	x = i / MaxSize
	y = i % MaxSize
	return x, y
}
