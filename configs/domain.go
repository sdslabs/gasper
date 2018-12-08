package configs

import (
	"encoding/json"
	"io/ioutil"
)

type domain struct {
	Name string `json:"domain"`
}

func getDomainName() string {
	file, _ := ioutil.ReadFile("../sample.json")
	var dom domain
	json.Unmarshal(file, &dom)
	return dom.Name
}
