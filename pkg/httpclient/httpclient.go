// Package httpclient contains functions for fetching data via http
package httpclient

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

// HTTPClient is the wrapper around http.Client for storing custom settings.
type HTTPClient struct {
	client *http.Client
	Header http.Header
}

// New creates new https client. Set `headers` as http headers. Use `timeout` as connection timeout.
func New(headers map[string]string, timeout time.Duration) *HTTPClient {
	client := &http.Client{Timeout: timeout}
	header := http.Header{}

	for k, v := range headers {
		header.Set(k, v)
	}

	return &HTTPClient{
		client: client,
		Header: header,
	}
}

// GetBody gets body data by url.
func (c *HTTPClient) GetBody(url string) (body []byte, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("httpclient: %v", err)
	}

	req.Header = c.Header

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

	return buf.Bytes(), nil
}
