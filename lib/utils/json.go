package utils

import (
	"encoding/json"
	"io/ioutil"
)

type redis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"DB"`
}

type service struct {
	Name   string `json:"name"`
	Deploy bool   `json:"deploy"`
	Port   string `json:"port"`
}

// JSONConfig is the type of parsed data from config.json
type JSONConfig struct {
	Domain   string    `json:"domain"`
	Redis    redis     `json:"redis"`
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
