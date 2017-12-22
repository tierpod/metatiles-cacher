package metatile

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

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

func decodeTile(r io.ReadSeeker, x, y int) ([]byte, error) {
	ml, err := decodeHeader(r)
	if err != nil {
		return nil, fmt.Errorf("decodeTile: %v", err)
	}

	size := int32(math.Sqrt(float64(ml.Count)))
	index := (int32(x)-ml.X)*size + (int32(y) - ml.Y)
	if index >= ml.Count {
		return nil, fmt.Errorf("decodeTile: invalid index %v/%v", index, ml.Count)
	}

	entry := ml.Index[index]
	if entry.Size > MaxEntrySize {
		return nil, fmt.Errorf("decodeTile: entry size > MaxEntrySize (size: %v)", entry.Size)
	}

	_, err = r.Seek(int64(entry.Offset), 0)
	if err != nil {
		return nil, fmt.Errorf("decodeTile: %v", err)
	}

	buf := make([]byte, entry.Size)
	l, err := r.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("decodeTile: %v", err)
	}

	if int32(l) != entry.Size {
		return nil, fmt.Errorf("decodeTile: invalid tile size: %v != %v", l, entry.Size)
	}

	return buf, nil
}
