// Package config contains functions for loading configuration file.
package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Service is the root of configuration.
type Service struct {
	Reader     ReaderSection     `yaml:"reader"`
	Writer     WriterSection     `yaml:"writer"`
	Zoom       ZoomSection       `yaml:"zoom"`
	Log        LogSection        `yaml:"log"`
	FileCache  FileCacheSection  `yaml:"filecache"`
	HTTPClient HTTPClientSection `yaml:"httpclient"`
	Sources    SourcesSection    `yaml:"sources"`
}

// ReaderSection contains Reader service configuration.
type ReaderSection struct {
	// Bind to address
	Bind string `yaml:"bind"`
	// Send requests to writer service?
	UseWriter bool `yaml:"use_writer"`
	// Send requests to remote sources?
	UseSources bool `yaml:"use_sources"`
	// writer service address
	WriterAddr string `yaml:"writer_addr"`
	// Token for XToken handler
	XToken string `yaml:"x_token"`
	// Cache-Control: max-age value in seconds
	MaxAge int `yaml:"max_age"`
}

// WriterSection contains Writer service configuration.
type WriterSection struct {
	Bind   string `yaml:"bind"`
	XToken string `yaml:"x_token"`
}

// ZoomSection contains min and max zoom levels.
type ZoomSection struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

// LogSection contains logger configuration.
type LogSection struct {
	Datetime bool `yaml:"datetime"`
	Debug    bool `yaml:"debug"`
}

// FileCacheSection contains file cache configuration.
type FileCacheSection struct {
	RootDir string `yaml:"root_dir"`
}

// HTTPClientSection contains http client configuration.
type HTTPClientSection struct {
	UserAgent string `yaml:"user_agent"`
}

// SourcesSection contains map of sources: key is name, value is url.
type SourcesSection struct {
	Map map[string]string
}

// UnmarshalYAML reads sources to Map["name"] = "url"
func (s *SourcesSection) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var url, source string
	var values yaml.MapSlice

	err := unmarshal(&values)
	if err != nil {
		return nil
	}

	m := make(map[string]string)

	for _, v := range values {
		url = v.Key.(string)
		source = v.Value.(string)
		m[url] = source
	}

	s.Map = m
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

	return &c, nil
}
