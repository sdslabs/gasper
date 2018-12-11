package utils

import (
	"encoding/json"
	"io/ioutil"
)

type service struct {
	Name   string   `json:"name"`
	Deploy bool     `json:"deploy"`
	Ports  []string `json:"ports"`
}

// JSONConfig is the type of parsed data from config.json
type JSONConfig struct {
	Domain   string    `json:"domain"`
	Services []service `json:"services"`
}

func parseJSON(path string) (JSONConfig, error) {
	var conf JSONConfig
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, err
	}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}

// SWSConfig is parsed data for `config.json`
var SWSConfig, _ = parseJSON("config.json")
