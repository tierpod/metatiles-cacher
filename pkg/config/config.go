// Package config contains functions for loading configuration file.
package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Service is the root of configuration.
type Service struct {
	Cacher     Cache      `yaml:"cacher"`
	Zoom       Zoom       `yaml:"zoom"`
	Log        Log        `yaml:"log"`
	FileCache  FileCache  `yaml:"filecache"`
	HTTPClient HTTPClient `yaml:"httpclient"`
	Sources    []Source   `yaml:"sources"`
	Test       string
}

// Source returns source configuration from Sources list by given name. If it does not exists,  returns error.
func (s Service) Source(name string) (Source, error) {
	for _, v := range s.Sources {
		if v.Name == name {
			return v, nil
		}
	}

	return Source{}, fmt.Errorf("source not found in sources")
}

// Cache contains cache service configuration.
type Cache struct {
	// Bind to address.
	Bind string `yaml:"bind"`
	// Send requests to remote source?
	UseSource bool `yaml:"use_source"`
	// Write to metatile cache?
	UseWriter bool `yaml:"use_writer"`
	// Token for XToken handler.
	XToken string `yaml:"x_token"`
	// Cache-Control: max-age value in seconds.
	MaxAge int `yaml:"max_age"`
}

// Zoom contains min and max zoom levels.
type Zoom struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

// Log contains logger configuration.
type Log struct {
	Datetime bool `yaml:"datetime"`
	Debug    bool `yaml:"debug"`
}

// FileCache contains file cache configuration.
type FileCache struct {
	RootDir string `yaml:"root_dir"`
}

// HTTPClient contains http client configuration.
type HTTPClient struct {
	UserAgent string `yaml:"user_agent"`
}

// Source contains source configuration.
type Source struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	CacheDir string `yaml:"cache_dir"`
	Region   Region `yaml:"region"`
}

// HasRegion return true if source has region. Otherwise return false.
func (s Source) HasRegion() bool {
	if s.Region.KML == "" && s.Region.Zoom.Min == 0 && s.Region.Zoom.Max == 0 {
		return false
	}

	return true
}

// Region contains region configuration.
type Region struct {
	KML      string `yaml:"kml"`
	Zoom     Zoom   `yaml:"zoom"`
	Polygons []byte
}

func (r *Region) readKML() error {
	data, err := ioutil.ReadFile(r.KML)
	if err != nil {
		return err
	}

	r.Polygons = data
	return nil
}

// Load loads yaml file and creates new service configuration.
func Load(path string) (*Service, error) {
	var c Service

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %v", err)
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config: %v", err)
	}

	// if source has "KML" section, read coordinates from given KML file
	for i := range c.Sources {
		if c.Sources[i].HasRegion() {
			err = c.Sources[i].Region.readKML()
			if err != nil {
				return nil, err
			}
		}
	}
	return &c, nil
}
