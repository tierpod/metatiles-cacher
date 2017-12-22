package metatile

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/tierpod/metatiles-cacher/pkg/tile"
)

// FetchEncodeTo fetchs tiles data, encodes it to metatile format and writes to w. Use ua as
// httpclient UserAgent.
func (mt Metatile) FetchEncodeTo(w io.Writer, url, ua string) error {
	var data Data
	data, err := mt.FetchTiles(url, ua)
	if err != nil {
		return err
	}

	err = mt.EncodeTiles(w, data)
	if err != nil {
		return err
	}
	return nil
}

// FetchTiles fetchs tiles data using url as template with placeholders {z}, {x}, {y} and return
// it. Use ua as httpclient UserAgent.
func (mt Metatile) FetchTiles(url, ua string) (data Data, err error) {
	xybox := mt.XYBox()
	for _, x := range xybox.X {
		for _, y := range xybox.Y {
			t := tile.New(mt.Zoom, x, y, "", mt.Map)
			offset := XYOffset(x, y)

			// fetch tile data to buffer
			var buf bytes.Buffer
			err := t.Fetch(&buf, url, ua)
			if err != nil {
				return Data{}, err
			}

			// read buffer to []byte
			d, err := ioutil.ReadAll(&buf)
			if err != nil {
				return Data{}, err
			}
			data[offset] = d
		}
	}

	// debug slow connections
	// time.Sleep(time.Second * 10)
	return data, nil
}
