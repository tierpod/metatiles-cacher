// Package metatile provides functions for decoding and encoding metatile files.
package metatile

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/tierpod/metatiles-cacher/pkg/coords"
)

const (
	// MaxCount is the maximum count of tiles in metatile.
	MaxCount = 1000
	// MaxEntrySize is the maximum size of metatile entry.
	MaxEntrySize = 1000000
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

// Encode encodes tiles to metatile and writes it to w.
func Encode(w io.Writer, m coords.Metatile, tiles [][]byte) error {
	// f.write(struct.pack("4s4i", META_MAGIC, METATILE * METATILE, x, y, z))
	x, y := m.MinXY()
	ml := &metaLayout{
		Magic: []byte{'M', 'E', 'T', 'A'},
		Count: int32(len(tiles)),
		X:     int32(x),
		Y:     int32(y),
		Z:     int32(m.Zoom),
	}
	// golang        |renderd.py
	// 20            |len(META_MAGIC) + 4 * 4
	// 8*len(tiles)  |(2 * 4) * (METATILE * METATILE)
	offset := int32(20 + 8*len(tiles))
	//size := t.ConvertToMeta().Size() // detect on zoom level?

	for i := 0; i < len(tiles); i++ {
		tile := tiles[i]
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

	if err := encodeHeader(w, ml); err != nil {
		return fmt.Errorf("metatile/encodeHeader: %v", err)
	}

	for i := 0; i < len(tiles); i++ {
		tile := tiles[i]

		if _, err := w.Write(tile); err != nil {
			return fmt.Errorf("metatile/write: %v", err)
		}
	}

	return nil
}

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
	if ml.Count > MaxCount {
		return nil, fmt.Errorf("Count > MaxCount (Count = %v)", ml.Count)
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

// GetTile decodes metatile from r and extract tile data.
func GetTile(r io.ReadSeeker, t coords.Tile) ([]byte, error) {
	ml, err := decodeHeader(r)
	if err != nil {
		return nil, err
	}

	size := int32(math.Sqrt(float64(ml.Count)))
	index := (int32(t.X)-ml.X)*size + (int32(t.Y) - ml.Y)
	if index >= ml.Count {
		return nil, fmt.Errorf("invalid index %v/%v", index, ml.Count)
	}

	entry := ml.Index[index]
	if entry.Size > MaxEntrySize {
		return nil, fmt.Errorf("entry size > MaxEntrySize (size: %v)", entry.Size)
	}

	_, err = r.Seek(int64(entry.Offset), 0)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, entry.Size)
	l, err := r.Read(buf)
	if err != nil {
		return nil, err
	}

	if int32(l) != entry.Size {
		return nil, fmt.Errorf("invalid tile size: %v != %v", l, entry.Size)
	}

	return buf, nil
}
