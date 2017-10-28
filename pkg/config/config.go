// Package config contains functions for loading configuration file.
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// Service is the root of configuration.
type Service struct {
	// metatiles_reader configuration
	Reader Reader `json:"reader"`
	// metatiles_writer configuration
	Writer Writer `json:"writer"`
	// sources for reader and writer
	Sources []Source `json:"sources"`
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

// Reader is the "reader" section of configuration.
type Reader struct {
	// Bind to address
	Bind string `json:"bind"`
	// Add datetime to log?
	LogDatetime bool `json:"log_datetime"`
	// Show debug messages in log?
	LogDebug bool `json:"log_debug"`
	// Max zoom level
	MaxZoom int `json:"max_zoom"`
	// Min zoom level
	MinZoom int `json:"min_zoom"`
	// Root directory for cache
	RootDir string `json:"root_dir"`
	// Writer service address. If "" - do not send request to writer.
	WriterAddr string `json:"writer_addr"`
	// Token for XToken handler
	XToken string `json:"x_token"`
}

// Writer is the "writer" section of configuration.
type Writer struct {
	Bind        string `json:"bind"`
	LogDatetime bool   `json:"log_datetime"`
	LogDebug    bool   `json:"log_debug"`
	MaxZoom     int    `json:"max_zoom"`
	MinZoom     int    `json:"min_zoom"`
	RootDir     string `json:"root_dir"`
	XToken      string `json:"x_token"`
}

// Source is the item in the "sources" section of configuration.
type Source struct {
	// Name of the source
	Name string `json:"name"`
	// Address of the source, contains "%v" will be replaced with "z/x/y.png"
	URL string `json:"url"`
}

// NewConfig loads json file and creates new service configuration.
func NewConfig(path string) *Service {
	var c Service

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &c)
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
