package configs

import (
	"encoding/json"
	"io/ioutil"
)

type domain struct {
	Name string `json:"domain"`
}

func getDomain() string {
	file, _ := ioutil.ReadFile("../config.json")
	var dom domain
	json.Unmarshal(file, &dom)
	return dom.Name
}
