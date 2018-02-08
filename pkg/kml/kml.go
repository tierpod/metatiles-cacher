// Package kml provides function for reading kml files from geofabrik.de and extracting Region
// information.
package kml

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/tierpod/go-osm/point"
)

// KML is the root space of kml file.
type KML struct {
	XMLName       xml.Name      `xml:"kml"`
	MultiGeometry MultiGeometry `xml:"Document>Placemark>MultiGeometry"`
}

// MultiGeometry is the Document>Placemark>MultiGeometry section of kml file. Contains Region.
type MultiGeometry struct {
	Region []Polygon `xml:"Polygon"`
}

// Polygon is the Document>Placemark>MultiGeometry>Polygon>outerBoundaryIs>LinearRing>coordinates section of kml file.
// Contains Coordinates.
type Polygon struct {
	Coordinates Coordinates `xml:"outerBoundaryIs>LinearRing>coordinates"`
}

// Coordinates is the slice of LatLong points.
type Coordinates []point.LatLong

// UnmarshalXML unmarshals string to LatLong coordinates.
func (c *Coordinates) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var content string
	if err := d.DecodeElement(&content, &start); err != nil {
		return err
	}

	var ll []point.LatLong

	for _, s := range strings.Split(strings.TrimSpace(content), "\n") {
		l := strings.Split(s, ",")

		long, err := strconv.ParseFloat(l[0], 64)
		if err != nil {
			return err
		}

		lat, err := strconv.ParseFloat(l[1], 64)
		if err != nil {
			return err
		}

		ll = append(ll, point.LatLong{Lat: lat, Long: long})
	}

	*c = ll
	return nil
}

// ExtractRegion extracts polygons of LinearRing points from kml file.
func ExtractRegion(r io.Reader) (point.Region, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	kml := KML{}
	err = xml.Unmarshal(b, &kml)
	if err != nil {
		return nil, err
	}

	var region point.Region
	for _, p := range kml.MultiGeometry.Region {
		polygon := point.Polygon(p.Coordinates)
		region = append(region, polygon)
	}

	return region, nil
}
