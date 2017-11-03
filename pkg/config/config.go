// Package config contains functions for loading configuration file.
package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Service is the root of configuration.
type Service struct {
	// metatiles_reader configuration
	ReaderSection ReaderSection `yaml:"reader"`
	// metatiles_writer configuration
	WriterSection WriterSection `yaml:"writer"`
	// sources for reader and writer
	ZoomSection       ZoomSection       `yaml:"zoom"`
	LogSection        LogSection        `yaml:"log"`
	FileCacheSection  FileCacheSection  `yaml:"filecache"`
	HTTPClientSection HTTPClientSection `yaml:"httpclient"`
	SourcesSection    SourcesSection    `yaml:"sources"`
	//Sources []Source `yaml:"sources"`
	// sources for reader and writer (in map)
	//SourcesMap map[string]string
}

// ReaderSection is the "reader" section of configuration.
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

// WriterSection is the "writer" section of configuration.
type WriterSection struct {
	Bind   string `yaml:"bind"`
	XToken string `yaml:"x_token"`
}

type ZoomSection struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

type LogSection struct {
	Datetime bool `yaml:"datetime"`
	Debug    bool `yaml:"debug"`
}

type FileCacheSection struct {
	RootDir string `yaml:"root_dir"`
}

type HTTPClientSection struct {
	UserAgent string `yaml:"user_agent"`
}

type SourcesSection struct {
	Sources map[string]string
}

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

	s.Sources = m
	return nil
}

// Load loads yaml file and creates new service configuration.
func Load(path string) *Service {
	var c Service

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Fatal(err)
	}

	return &c
}
