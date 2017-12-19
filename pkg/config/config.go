// Package config contains functions for loading configuration file.
package config

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/tierpod/metatiles-cacher/pkg/polygon"

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
	Service   Service   `yaml:"service"`
	Log       Log       `yaml:"log"`
	FileCache FileCache `yaml:"filecache"`
	Fetch     Fetch     `yaml:"fetch"`
	Sources   []Source  `yaml:"sources"`
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

// Fetch contains fetcher configuration.
type Fetch struct {
	UserAgent    string `yaml:"user_agent"`
	QueueTimeout int    `yaml:"queue_timeout"`
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
	Zoom     Zoom   `yaml:"zoom"`
	Polygons polygon.Region
}

func (r *Region) readFile() error {
	var region polygon.Region
	var err error

	switch path.Ext(r.File) {
	case ".yaml", ".yml":
		region, err = readYAML(r.File)
	case ".kml":
		region, err = readKML(r.File)
	default:
		return fmt.Errorf("readFile: unknown file format: %v", path.Ext(r.File))
	}

	if err != nil {
		return fmt.Errorf("readFile: %v", err)
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

	if c.Fetch.QueueTimeout == 0 {
		c.Fetch.QueueTimeout = 30
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

		if c.Sources[i].HasRegion() {
			// if Source.Region has "File" section read coordinates from given file to Region.Polygons struct
			err = c.Sources[i].Region.readFile()
			if err != nil {
				return nil, err
			}
			// if Source.Region.Zoom is not set, use defaults.
			if c.Sources[i].Region.Zoom.Min == 0 && c.Sources[i].Region.Zoom.Max == 0 {
				c.Sources[i].Region.Zoom.Min = DefaultMinZoom
				c.Sources[i].Region.Zoom.Max = DefaultMaxZoom
			}
		}
	}
	return &c, nil
}
