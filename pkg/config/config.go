// Package config contains functions for loading configuration file.
package config

import (
	"fmt"
	"io/ioutil"

	"github.com/tierpod/metatiles-cacher/pkg/coords"

	"gopkg.in/yaml.v2"
)

const (
	// DefaultMinZoom is the default minimum zoom level.
	DefaultMinZoom = 1
	// DefaultMaxZoom is the default maximum zoom level.
	DefaultMaxZoom = 18
)

// Config is the root of configuration.
type Config struct {
	Service    Service    `yaml:"service"`
	Log        Log        `yaml:"log"`
	FileCache  FileCache  `yaml:"filecache"`
	HTTPClient HTTPClient `yaml:"httpclient"`
	Sources    []Source   `yaml:"sources"`
}

// Source returns source configuration from Sources list by given name. If it does not exist, return error.
func (c Config) Source(name string) (Source, error) {
	for _, v := range c.Sources {
		if v.Name == name {
			return v, nil
		}
	}

	return Source{}, fmt.Errorf("source not found in sources")
}

// Service contains metatiles-cacher service configuration.
type Service struct {
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
	Zoom     Zoom   `yaml:"zoom"`
	Region   Region `yaml:"region"`
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
	Polygons []coords.Polygon
}

func (r *Region) readFile() error {
	type yamlPolygon struct {
		Polygon coords.Polygon `yaml:"polygon"`
	}

	var result []yamlPolygon

	data, err := ioutil.ReadFile(r.File)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &result)

	// convert yaml struct to coords struct
	var p []coords.Polygon
	for _, v := range result {
		p = append(p, v.Polygon)
	}

	r.Polygons = p
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

	for i := range c.Sources {
		// if Source.Zoom is not set, use defaults.
		if c.Sources[i].Zoom.Min == 0 && c.Sources[i].Zoom.Max == 0 {
			c.Sources[i].Zoom.Min = DefaultMinZoom
			c.Sources[i].Zoom.Max = DefaultMaxZoom
		}

		// if Source.CacheDir is not set, use Source.Name.
		if c.Sources[i].CacheDir == "" {
			c.Sources[i].CacheDir = c.Sources[i].Name
		}

		// if Source.Region has "File" section, read coordinates from given file to Region.Polygons struct.
		if c.Sources[i].HasRegion() {
			err = c.Sources[i].Region.readFile()
			if err != nil {
				return nil, err
			}
		}
	}
	return &c, nil
}
