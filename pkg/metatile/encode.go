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

// MetaLayout is the metatile file struct.
type MetaLayout struct {
	Magic   []byte
	Count   int32
	X, Y, Z int32
	Index   []metaEntry
}

func (m *MetaLayout) String() string {
	return fmt.Sprintf("MetatileLayout{X:%v Y:%v Z:%v Count:%v}", m.X, m.Y, m.Z, m.Count)
}

// EncodeHeader encodes ml and writes it to w.
func EncodeHeader(w io.Writer, ml *MetaLayout) error {
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

// EncodeData encodes data for metatile mt and writes it to w.
func (mt Metatile) EncodeData(w io.Writer, data Data) error {
	mSize := MaxSize * MaxSize

	if len(data) < Area {
		return fmt.Errorf("encodeMetatile: data size: %v < %v", len(data), mSize)
	}

	ml := &MetaLayout{
		Magic: []byte{'M', 'E', 'T', 'A'},
		Count: int32(mSize),
		X:     int32(mt.X),
		Y:     int32(mt.Y),
		Z:     int32(mt.Zoom),
	}
	offset := int32(20 + 8*mSize)

	// calculate offsets and sizes
	for i := 0; i < mSize; i++ {
		tile := data[i]
		s := int32(len(tile))
		if s > MaxEntrySize {
			return fmt.Errorf("encodeMetatile: entry size > MaxEntrySize (size: %v)", s)
		}

		ml.Index = append(ml.Index, metaEntry{
			Offset: offset,
			Size:   s,
		})
		offset += s
	}

	// fmt.Printf("%+v\n", ml)

	// encode and write headers
	if err := EncodeHeader(w, ml); err != nil {
		return fmt.Errorf("encodeMetatile: %v", err)
	}

	// encode and write data
	for i := 0; i < len(data); i++ {
		tile := data[i]

		if _, err := w.Write(tile); err != nil {
			return fmt.Errorf("encodeMetatile: %v", err)
		}
	}

	return nil
}

// EncodeTiles encodes tiles data to metatile format for mt and writes it to w.
func (mt Metatile) EncodeTiles(w io.Writer, data Data) error {
	err := mt.EncodeData(w, data)
	if err != nil {
		return err
	}

	return nil
}
