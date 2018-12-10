package utils

import (
	"encoding/json"
	"io/ioutil"
)

func parseJSON(location string) (map[string]interface{}, error) {
	file, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	var parsedData map[string]interface{}
	err = json.Unmarshal(file, &parsedData)
	if err != nil {
		return nil, err
	}
	return parsedData, nil
}

// SWSConfig is parsed data for `config.json`
var SWSConfig, _ = parseJSON("config.json")

// Shortcuts - start with `Config`
var (
	ConfigDomain = SWSConfig["domain"].(string)
)
