package config

import (
	"fmt"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	// file not found
	_, err := Load("testdata/notfound")
	if err.Error() != "read config: open testdata/notfound: no such file or directory" {
		t.Errorf("Load: expected \"no such file\" error, got %v", err)
	}

	_, err = Load("testdata/config.yaml")
	if err != nil {
		t.Errorf("Load: expected no error, got %v", err)
	}
}

func TestSource(t *testing.T) {
	testSource := Source{
		URL:     "http://tilesrv1/style/{z}/{x}/{y}.png",
		MaxZoom: 18,
	}

	config, _ := Load("testdata/config.yaml")
	// source not found
	_, err := config.Source("notfound")
	if err == nil {
		t.Errorf("Source: expected \"source not found\" error, got nil")
	}

	// source found
	source, err := config.Source("testsrc1")
	if err != nil {
		t.Errorf("Source: expected no error, got %v", err)
	}
	if !reflect.DeepEqual(testSource, source) {
		t.Errorf("Source: testSource(%+v) != source(%+v)", testSource, source)
	}
}

func ExampleLoad() {
	c, _ := Load("testdata/config.yaml")
	s, err := c.Source("testsrc2")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("source.URL: '%v'\n", s.URL)
	fmt.Printf("source.MaxZoom: %+v\n", s.MaxZoom)
	fmt.Printf("source.LastUpdate: '%v'\n", s.LastUpdate)
	fmt.Printf("source.Region.File: '%v'\n", s.Region.File)
	fmt.Printf("source.Region.MaxZoom: %+v\n", s.Region.MaxZoom)

	// Output:
	// source.URL: 'http://testsrv2/style/{z}/{x}/{y}.png?key=123'
	// source.MaxZoom: 14
	// source.LastUpdate: '0001-01-01 00:00:00 +0000 UTC'
	// source.Region.File: 'testdata/test_region.kml'
	// source.Region.MaxZoom: 17
}
