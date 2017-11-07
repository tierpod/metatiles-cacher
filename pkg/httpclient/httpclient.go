// Package httpclient contains functions for send http requests
package httpclient

import (
	"fmt"
	"io"
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

// PostJSON posts json data by url
func PostJSON(url string, body io.Reader) error {
	res, err := http.Post(url, "application/json; charset=utf-8", body)
	if err != nil {
		return fmt.Errorf("httpclient.PostJSON: %v: %v", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("httpclient.PostJSON: %v: Response status %v", url, res.StatusCode)
	}

	return nil
}
