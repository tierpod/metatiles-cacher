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

	// unmarshal error
	_, err = Load("testdata/config2.yaml")
	if err.Error() != "unmarshal config: yaml: line 3: could not find expected ':'" {
		t.Errorf("Load: expected \"unmarshal error\" error, got %v", err)
	}

	// unknown region format
	_, err = Load("testdata/config3.yaml")
	if err != nil && err.Error() != "readFile: unknown file format: .unknown" {
		t.Errorf("Load: expected \"unknown format\" error, got %v", err)
	}

	_, err = Load("testdata/config.yaml")
	if err != nil {
		t.Errorf("Load: expected no error, got %v", err)
	}
	//fmt.Printf("%+v\n", config)
}

func TestSource(t *testing.T) {
	testSource := Source{
		Name:     "testsrc1",
		URL:      "http://tilesrv1/style/{tile}",
		CacheDir: "testsrc1",
		Zoom: Zoom{
			Min: 1,
			Max: 18,
		},
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
	config, _ := Load("testdata/config.yaml")
	for _, s := range config.Sources {
		fmt.Printf("source.Name: '%v'\n", s.Name)
		fmt.Printf("source.URL: '%v'\n", s.URL)
		fmt.Printf("source.CacheDir: '%v'\n", s.CacheDir)
		fmt.Printf("source.Zoom: %+v\n", s.Zoom)
		fmt.Printf("source.Region.File: '%v'\n", s.Region.File)
		fmt.Printf("source.Region.Zoom: %+v\n", s.Region.Zoom)
		fmt.Println("---")
	}

	// Output:
	// source.Name: 'testsrc1'
	// source.URL: 'http://tilesrv1/style/{tile}'
	// source.CacheDir: 'testsrc1'
	// source.Zoom: {Min:1 Max:18}
	// source.Region.File: ''
	// source.Region.Zoom: {Min:0 Max:0}
	// ---
	// source.Name: 'testsrc2'
	// source.URL: 'http://testsrv2/style/{tile}?api_key=123'
	// source.CacheDir: 'test'
	// source.Zoom: {Min:1 Max:18}
	// source.Region.File: ''
	// source.Region.Zoom: {Min:0 Max:0}
	// ---
	// source.Name: 'testsrc3'
	// source.URL: 'http://testsrv3/style/{tile}'
	// source.CacheDir: 'test'
	// source.Zoom: {Min:1 Max:19}
	// source.Region.File: 'testdata/test_region.yaml'
	// source.Region.Zoom: {Min:1 Max:10}
	// ---
	// source.Name: 'testsrc4'
	// source.URL: 'http://testsrv4/style/{tile}'
	// source.CacheDir: 'test'
	// source.Zoom: {Min:1 Max:19}
	// source.Region.File: 'testdata/test_region.kml'
	// source.Region.Zoom: {Min:1 Max:18}
	// ---
}
