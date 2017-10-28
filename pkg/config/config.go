// Package config contains functions for loading configuration file.
package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Service is the root of configuration.
type Service struct {
	// metatiles_reader configuration
	Reader ReaderSection `yaml:"reader"`
	// metatiles_writer configuration
	Writer WriterSection `yaml:"writer"`
	// sources for reader and writer
	Sources []Source `yaml:"sources"`
	// sources for reader and writer (in map)
	SourcesMap map[string]string
}

// SourceInfo returns information about source from config: {name: url}. Return err if source not found.
func (s Service) SourceInfo(source string) (name string, url string, err error) {
	for _, v := range s.Sources {
		if v.Name == source {
			return v.Name, v.URL, nil
		}
	}

	return "", "", fmt.Errorf("source for style %v not found", source)
}

// ReaderSection is the "reader" section of configuration.
type ReaderSection struct {
	// Bind to address
	Bind string `yaml:"bind"`
	// Add datetime to log?
	LogDatetime bool `yaml:"log_datetime"`
	// Show debug messages in log?
	LogDebug bool `yaml:"log_debug"`
	// Max zoom level
	MaxZoom int `yaml:"max_zoom"`
	// Min zoom level
	MinZoom int `yaml:"min_zoom"`
	// Root directory for cache
	RootDir string `yaml:"root_dir"`
	// Writer service address. If "" - do not send request to writer.
	WriterAddr string `yaml:"writer_addr"`
	// Token for XToken handler
	XToken string `yaml:"x_token"`
}

// WriterSection is the "writer" section of configuration.
type WriterSection struct {
	Bind        string `yaml:"bind"`
	LogDatetime bool   `yaml:"log_datetime"`
	LogDebug    bool   `yaml:"log_debug"`
	MaxZoom     int    `yaml:"max_zoom"`
	MinZoom     int    `yaml:"min_zoom"`
	RootDir     string `yaml:"root_dir"`
	XToken      string `yaml:"x_token"`
}

// Source is the item in the "sources" section of configuration, contains name and url.
type Source struct {
	// Name of the source
	Name string `yaml:"name"`
	// Address of the source, contains "%v" will be replaced with "z/x/y.png"
	URL string `yaml:"url"`
}

// NewConfig loads yaml file and creates new service configuration.
func NewConfig(path string) *Service {
	var c Service

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Fatal(err)
	}

	c.SourcesMap = c.sourcesToMap()

	return &c
}

func (s Service) sourcesToMap() map[string]string {
	result := make(map[string]string)
	for _, v := range s.Sources {
		result[v.Name] = v.URL
	}

	return result
}
