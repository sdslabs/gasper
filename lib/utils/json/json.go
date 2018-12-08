package json

import (
	encodingJSON "encoding/json"
	"io/ioutil"
)

// ParseJSONFromFile takes location of the file and returns map[string]interface{}
// which is a way of representing unstructured json data
func ParseJSONFromFile(location string) (map[string]interface{}, error) {
	file, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	var parsedData map[string]interface{}
	err = encodingJSON.Unmarshal(file, &parsedData)
	if err != nil {
		return nil, err
	}
	return parsedData, nil
}

// SWSConfig is the configuration of SWS server
var SWSConfig, _ = ParseJSONFromFile("../../config.json")
