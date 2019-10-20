package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func parseJSON(path string) *GasperCfg {
	config := &GasperCfg{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("File %s does not exist", path))
	}
	err = json.Unmarshal(file, config)
	if err != nil {
		panic(fmt.Sprintf("Invalid %s file", path))
	}
	return config
}
