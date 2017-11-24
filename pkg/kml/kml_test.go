package kml

import (
	"fmt"
	"os"
	"testing"
)

func TestExtractRegion(t *testing.T) {
	var err error
	// extract longitude
	file, _ := os.Open("testdata/test_region2.kml")
	defer file.Close()

	_, err = ExtractRegion(file)
	if err.Error() != "strconv.ParseFloat: parsing \"wrong\": invalid syntax" {
		t.Errorf("ExtractRegion: expected strconv.ParseFloat error")
	}

	// extract latitude
	file, _ = os.Open("testdata/test_region3.kml")
	defer file.Close()

	_, err = ExtractRegion(file)
	if err.Error() != "strconv.ParseFloat: parsing \"wrong\": invalid syntax" {
		t.Errorf("ExtractRegion: expected strconv.ParseFloat error")
	}
}

func ExampleExtractRegion() {
	file, _ := os.Open("testdata/test_region.kml")
	defer file.Close()

	region, _ := ExtractRegion(file)
	fmt.Println("region:")
	for _, v := range region {
		fmt.Println(" polygon:", v)
	}

	// Output:
	// region:
	//  polygon: [LatLong{77.7--168.25} LatLong{77.7--179.9999} LatLong{58.1--179.9999} LatLong{58.1--168.25} LatLong{77.7--168.25}]
	//  polygon: [LatLong{84.52666-39.6908} LatLong{84.38487-179.9999} LatLong{26.27883-179.9999} LatLong{22.062707-142.084541} LatLong{84.52666-39.6908}]
}
