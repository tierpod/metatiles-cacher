// Package metatile provides functions for decoding and encoding metatile files.
//
// Metatile format description: https://github.com/openstreetmap/mod_tile/blob/master/src/metatile.cpp
package metatile

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	// MaxCount is the maximum count of tiles in metatile.
	MaxCount = 1000
	// MaxEntrySize is the maximum size of metatile entry in bytes.
	MaxEntrySize = 2000000
)

type metaEntry struct {
	Offset int32
	Size   int32
}

type metaLayout struct {
	Magic   []byte
	Count   int32
	X, Y, Z int32
	Index   []metaEntry
}

func encodeHeader(w io.Writer, ml *metaLayout) error {
	endian := binary.LittleEndian
	var err error
	if err = binary.Write(w, endian, ml.Magic); err != nil {
		return err
	}
	if err = binary.Write(w, endian, ml.Count); err != nil {
		return err
	}
	if err = binary.Write(w, endian, ml.X); err != nil {
		return err
	}
	if err = binary.Write(w, endian, ml.Y); err != nil {
		return err
	}
	if err = binary.Write(w, endian, ml.Z); err != nil {
		return err
	}
	for _, ent := range ml.Index {
		if err = binary.Write(w, endian, ent); err != nil {
			return err
		}
	}
	return nil
}

// Encode encodes data to metatile and writes it to w.
func (m Metatile) Encode(w io.Writer, data Data) error {
	mSize := MaxSize * MaxSize

	if len(data) < Area {
		return fmt.Errorf("data size: %v < %v", len(data), mSize)
	}

	ml := &metaLayout{
		Magic: []byte{'M', 'E', 'T', 'A'},
		Count: int32(mSize),
		X:     int32(m.X),
		Y:     int32(m.Y),
		Z:     int32(m.Zoom),
	}
	offset := int32(20 + 8*mSize)

	// calculate offsets and sizes
	for i := 0; i < mSize; i++ {
		tile := data[i]
		s := int32(len(tile))
		if s > MaxEntrySize {
			return fmt.Errorf("entry size > MaxEntrySize (size: %v)", s)
		}

		ml.Index = append(ml.Index, metaEntry{
			Offset: offset,
			Size:   s,
		})
		offset += s
	}

	// fmt.Printf("%+v\n", ml)

	// encode and write headers
	if err := encodeHeader(w, ml); err != nil {
		return fmt.Errorf("metatile/encodeHeader: %v", err)
	}

	// encode and write data
	for i := 0; i < len(data); i++ {
		tile := data[i]

		if _, err := w.Write(tile); err != nil {
			return fmt.Errorf("metatile/write: %v", err)
		}
	}

	return nil
}
