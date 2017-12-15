// Package httpclient contains functions for fetching data via http
package httpclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
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
