package fetch

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tierpod/metatiles-cacher/pkg/httpclient"
	"github.com/tierpod/metatiles-cacher/pkg/tile"
)

// Tile fetchs tile data, using URLTmpl as URL template with placeholders {x} {y} {z}.
func (f *Fetch) Tile(t tile.Tile, URLTmpl string) (tile.Data, error) {
	url := strings.Replace(URLTmpl, "{z}", strconv.Itoa(t.Zoom), 1)
	url = strings.Replace(url, "{x}", strconv.Itoa(t.X), 1)
	url = strings.Replace(url, "{y}", strconv.Itoa(t.Y), 1)

	f.logger.Printf("Fetch/Tile: get from URL(%v)", url)

	data, err := httpclient.Get(url, f.cfg.UserAgent)
	if err != nil {
		f.logger.Printf("[ERROR] Fetch/Tile: %v", err)
		return nil, fmt.Errorf("Fetch/Tile: %v", err)
	}

	return data, nil
}
