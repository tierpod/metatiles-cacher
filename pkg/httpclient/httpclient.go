// Package httpclient contains functions for fetching data via http
package httpclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/tierpod/metatiles-cacher/pkg/metatile"
)

// Get gets data by url
func Get(url, ua string) (data []byte, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("httpclient/Newrequest: %v", err)
	}

	req.Header.Set("User-Agent", ua)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpclient/Get: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("httpclient/Get: %v: Response status %v", url, resp.StatusCode)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("httpclient/Get: %v", err)
	}
	return data, nil
}

// FetchMetatile fetchs metatile data from xybox coordinates for given metatile.
func FetchMetatile(m metatile.Metatile, ext, sURL, ua string) (metatile.Data, error) {
	var data metatile.Data
	xybox := m.XYBox()

	for _, x := range xybox.X {
		for _, y := range xybox.Y {
			offset := metatile.XYOffset(x, y)
			tile := strconv.Itoa(m.Zoom) + "/" + strconv.Itoa(x) + "/" + strconv.Itoa(y) + ext
			url := strings.Replace(sURL, "{tile}", tile, 1)
			res, err := Get(url, ua)
			if err != nil {
				return data, err
			}
			data[offset] = res
		}
	}

	return data, nil
}
