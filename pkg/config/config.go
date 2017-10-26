// Package config contains functions for loading service config.
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// Service is the root of configuration.
type Service struct {
	Reader  Reader   `json:"reader"`
	Writer  Writer   `json:"writer"`
	Sources []Source `json:"sources"`
}

// GetStyleInfo returns information about style from config: name, url
func (s Service) GetStyleInfo(style string) (name string, urlTmpl string, err error) {
	for _, v := range s.Sources {
		if v.Name == style {
			return v.Name, v.URL, nil
		}
	}

	return "", "", fmt.Errorf("source for style %v not found", style)
}

// SourcesToMap converts Sources to map: name=url
func (s Service) SourcesToMap() map[string]string {
	result := make(map[string]string)
	for _, v := range s.Sources {
		result[v.Name] = v.URL
	}

	return result
}

// Reader is the "reader" section of configuration
type Reader struct {
	Bind        string `json:"bind"`
	LogDebug    bool   `json:"log_debug"`
	LogDatetime bool   `json:"log_datetime"`
	MinZoom     int    `json:"min_zoom"`
	MaxZoom     int    `json:"max_zoom"`
	XToken      string `json:"x_token"`
	RootDir     string `json:"root_dir"`
	Writer      bool   `json:"writer"`
	WriterAddr  string `json:"writer_addr"`
}

// Writer is the "writer" section of configuration
type Writer struct {
	Bind        string `json:"bind"`
	LogDebug    bool   `json:"log_debug"`
	LogDatetime bool   `json:"log_datetime"`
	XToken      string `json:"x_token"`
	RootDir     string `json:"root_dir"`
}

// Source is the item in the "sources" section
type Source struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// NewConfig creates new ServiceConf from file
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

	return &c
}
