// Package httpclient contains functions for fetching data via http
package httpclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/tierpod/metatiles-cacher/pkg/coords"
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

// FetchMetatileData fetchs metatile data inside xybox coordinates for given zoom and ext.
func FetchMetatileData(xybox coords.XYBox, zoom int, ext, sURL, ua string) (coords.MetatileData, error) {
	var data coords.MetatileData

	for _, x := range xybox.X {
		for _, y := range xybox.Y {
			mo := coords.XYToMetatileOffset(x, y)
			tile := strconv.Itoa(zoom) + "/" + strconv.Itoa(x) + "/" + strconv.Itoa(y) + `.` + ext
			url := strings.Replace(sURL, "{tile}", tile, 1)
			res, err := Get(url, ua)
			if err != nil {
				return data, err
			}
			data[mo] = res
		}
	}

	return data, nil
}
