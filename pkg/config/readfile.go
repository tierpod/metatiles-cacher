package config

import (
	"io/ioutil"
	"os"

	"github.com/tierpod/metatiles-cacher/pkg/coords"
	"github.com/tierpod/metatiles-cacher/pkg/kml"
	"gopkg.in/yaml.v2"
)

// readYAML reads yaml file and convert it to coords.Region.
func readYAML(path string) (coords.Region, error) {
	type yamlPolygon struct {
		Polygon coords.Polygon `yaml:"polygon"`
	}

	var result []yamlPolygon

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	// convert yaml struct to coords struct
	var region coords.Region
	for _, v := range result {
		region = append(region, v.Polygon)
	}

	return region, nil
}

// readKML reads kml file from geofabric.de and convert it to coords.Region.
func readKML(path string) (coords.Region, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	region, err := kml.ExtractRegion(file)
	if err != nil {
		return nil, err
	}

	return region, nil
}
