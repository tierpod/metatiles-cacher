// Package config contains functions for loading configuration file.
package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/tierpod/go-osm/point"
	"github.com/tierpod/metatiles-cacher/pkg/kml"
	"gopkg.in/yaml.v2"
)

const (
	// MinZoom is the default source min_zoom value.
	MinZoom = 1
	// MaxZoom is the default source max_zoom value.
	MaxZoom = 18
	// MediumZoom used in cache expire
	MediumZoom = 13
	// LowZoom used in cache expire
	LowZoom = 9
)

// Config is the root of configuration.
type Config struct {
	Cache      Cache             `yaml:"cache"`
	Fetch      Fetch             `yaml:"fetch"`
	HTTP       HTTP              `yaml:"http"`
	HTTPClient HTTPClient        `yaml:"httpclient"`
	Log        Log               `yaml:"log"`
	Sources    map[string]Source `yaml:"sources"`
}

// Source returns source configuration from Sources list by given name. If it does not exist, return error.
func (c Config) Source(name string) (Source, error) {
	if s, found := c.Sources[name]; found {
		return s, nil
	}

	return Source{}, fmt.Errorf("source not found in sources")
}

// HTTP contains web service configuration.
type HTTP struct {
	// Bind to address.
	Bind string `yaml:"bind"`
	// Send requests to remote source?
	XToken string `yaml:"x_token"`
	// HTTP server headers
	Headers map[string]string `yaml:"headers"`
}

// HTTPClient contains http client configuration.
type HTTPClient struct {
	Headers map[string]string `yaml:"headers"`
	Timeout time.Duration     `yaml:"timeout"`
}

// Log contains logger configuration.
type Log struct {
	Datetime bool `yaml:"datetime"`
	Debug    bool `yaml:"debug"`
}

// Cache contains file cache configuration.
type Cache struct {
	Dir string `yaml:"dir"`
}

// Fetch contains fetchsvc.Service configuration
type Fetch struct {
	Enabled bool `yaml:"enabled"`
	Workers int  `yaml:"workers"`
	Buffer  int  `yaml:"buffer"`
}

// Source contains source configuration.
type Source struct {
	URL     string `yaml:"url"`
	MaxZoom int    `yaml:"max_zoom"`
	Region  Region `yaml:"region"`
}

// HasRegion return true if source has region section. Otherwise return false.
func (s Source) HasRegion() bool {
	if s.Region.File == "" {
		return false
	}

	return true
}

// Region contains region configuration.
type Region struct {
	File     string `yaml:"file"`
	MaxZoom  int    `yaml:"max_zoom"`
	Polygons point.Region
}

func (r *Region) readKML(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	region, err := kml.ExtractRegion(file)
	if err != nil {
		return err
	}

	r.Polygons = region
	return nil
}

// Load loads yaml file and creates new service configuration.
func Load(path string) (*Config, error) {
	var c Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %v", err)
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config: %v", err)
	}

	// convert timeout to seconds
	c.HTTPClient.Timeout = c.HTTPClient.Timeout * time.Second

	for name, s := range c.Sources {
		// if Source.MaxZoom is not set, use defaults.
		if s.MaxZoom == 0 {
			s.MaxZoom = MaxZoom
		}

		if s.HasRegion() {
			// if Source.Region has "File" section read coordinates from given file to Region.Polygons struct
			err = s.Region.readKML(s.Region.File)
			if err != nil {
				return nil, err
			}

			// if Source.Region.Zoom is not set, use defaults.
			if s.Region.MaxZoom == 0 {
				s.Region.MaxZoom = s.MaxZoom
			}
		}

		c.Sources[name] = s
	}

	return &c, nil
}
