// Package httpclient contains functions for fetching data via http
package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// GetBody gets body data by url.
func GetBody(url, ua string) (body io.Reader, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("httpclient: %v", err)
	}

	req.Header.Set("User-Agent", ua)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpclient: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("httpclient: invalid response status code %v for %v", url, resp.StatusCode)
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("httpclient: %v", err)
	}

	return &buf, nil
}
