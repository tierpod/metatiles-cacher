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
	//  polygon: [{Lat:77.7 Long:-168.25} {Lat:77.7 Long:-179.9999} {Lat:58.1 Long:-179.9999} {Lat:58.1 Long:-168.25} {Lat:77.7 Long:-168.25}]
	//  polygon: [{Lat:84.52666 Long:39.6908} {Lat:84.38487 Long:179.9999} {Lat:26.27883 Long:179.9999} {Lat:22.062707 Long:142.084541} {Lat:84.52666 Long:39.6908}]
}
