package config

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetStyleInfo(t *testing.T) {
	source1 := Source{Name: "style1", URL: "http://test1/style1/%v"}
	source2 := Source{Name: "style2", URL: "http://test2/style2/%v"}
	service := Service{Sources: []Source{source1, source2}}

	validName := "style1"
	validURL := "http://test1/style1/%v"
	name, url, err := service.GetStyleInfo("style1")
	if err != nil {
		t.Errorf("GetStyleInfo: got error %v", err)
	}
	if name != validName || url != validURL {
		t.Errorf("GetStyleInfo: expected {name:%v style:%v}, got: {name:%v style:%v}", validName, validURL, name, url)
	}

	invalidName := "notexist"
	invalidErr := fmt.Sprintf("source for style %v not found", invalidName)
	_, _, err = service.GetStyleInfo(invalidName)
	if err.Error() != invalidErr {
		t.Errorf("GetStyleInfo: expected err \"%v\", got \"%v\"", invalidErr, err.Error())
	}
}

func TestSourcesToMap(t *testing.T) {
	source1 := Source{Name: "style1", URL: "http://test1/style1/%v"}
	source2 := Source{Name: "style2", URL: "http://test2/style2/%v"}
	service := Service{Sources: []Source{source1, source2}}

	result := make(map[string]string)
	result["style1"] = "http://test1/style1/%v"
	result["style2"] = "http://test2/style2/%v"

	got := service.SourcesToMap()
	if !reflect.DeepEqual(got, result) {
		t.Errorf("SourcesToMap: expected %v, got %v", result, got)
	}
}
