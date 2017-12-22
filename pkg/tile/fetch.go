package tile

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/tierpod/metatiles-cacher/pkg/httpclient"
)

// Fetch fetchs tile data from url with placeholders {z}, {x}, {y} and writes it to w. Use ua as
// httpclient UserAgent.
func (t Tile) Fetch(w io.Writer, url, ua string) error {
	url = strings.Replace(url, "{z}", strconv.Itoa(t.Zoom), 1)
	url = strings.Replace(url, "{x}", strconv.Itoa(t.X), 1)
	url = strings.Replace(url, "{y}", strconv.Itoa(t.Y), 1)

	body, err := httpclient.GetBody(url, ua)
	if err != nil {
		return fmt.Errorf("GetBody: %v", err)
	}

	_, err = io.Copy(w, body)
	if err != nil {
		return fmt.Errorf("io.copy: %v", err)
	}
	return nil
}
