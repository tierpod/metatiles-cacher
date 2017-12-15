package config

import (
	"io/ioutil"
	"os"

	"github.com/tierpod/metatiles-cacher/pkg/kml"
	"github.com/tierpod/metatiles-cacher/pkg/polygon"
	"gopkg.in/yaml.v2"
)

// readYAML reads yaml file and convert it to coords.Region.
func readYAML(path string) (polygon.Region, error) {
	type yamlPolygon struct {
		Polygon polygon.Polygon `yaml:"polygon"`
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
	var region polygon.Region
	for _, v := range result {
		region = append(region, v.Polygon)
	}

	return region, nil
}

// readKML reads kml file from geofabric.de and convert it to coords.Region.
func readKML(path string) (polygon.Region, error) {
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
