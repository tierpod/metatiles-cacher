// Package kml provides function for reading kml files from geofabrik.de
package kml

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/tierpod/metatiles-cacher/pkg/coords"
)

// KML is the root space of kml file.
type KML struct {
	XMLName       xml.Name      `xml:"kml"`
	MultiGeometry MultiGeometry `xml:"Document>Placemark>MultiGeometry"`
}

// MultiGeometry is the Document>Placemark>MultiGeometry section of kml file. Contains Polygons.
type MultiGeometry struct {
	Polygons []Polygon `xml:"Polygon"`
}

// Polygon is the Document>Placemark>MultiGeometry>Polygon>outerBoundaryIs>LinearRing>coordinates section of kml file.
// Contains Coordinates.
type Polygon struct {
	//Coordinates string `xml:"outerBoundaryIs>LinearRing>coordinates"`
	Coordinates Coordinates `xml:"outerBoundaryIs>LinearRing>coordinates"`
}

// Coordinates is the slice of LatLong points.
type Coordinates []coords.LatLong

// UnmarshalXML unmarshals string to LatLong coordinates.
func (c *Coordinates) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var content string
	if err := d.DecodeElement(&content, &start); err != nil {
		return err
	}

	var ll []coords.LatLong

	for _, s := range strings.Split(strings.TrimSpace(content), "\n") {
		l := strings.Split(s, ",")
		lat, err := strconv.ParseFloat(l[0], 64)
		if err != nil {
			return err
		}

		long, err := strconv.ParseFloat(l[1], 64)
		if err != nil {
			return err
		}

		ll = append(ll, coords.LatLong{Lat: lat, Long: long})
	}

	*c = ll
	return nil
}

// ExtractPolygons extracts polygons of LinearRing points from kml file.
func ExtractPolygons(r io.Reader) (coords.Region, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	kml := KML{}
	err = xml.Unmarshal(b, &kml)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("ep: %+v\n", kml)

	//var region coords.Region
	//var polygons []coords.Polygon
	//for _, p := range kml.MultiGeometry.Polygons {
	//fmt.Println(k, v)
	//for _, c := range strings.Split(p.Coordinates, "\n") {
	//fmt.Println(kk, s)
	//polygons = apend(polygons, c)

	//}
	//region = append(region, polygons)
	//}

	//var region coords.Region
	for _, p := range kml.MultiGeometry.Polygons {
		fmt.Printf("polygon: %+v\n", p)
	}

	return nil, nil
}
