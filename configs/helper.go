package configs

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

func parseTOML(path string) *GasperCfg {
	config := &GasperCfg{}
	if _, err := toml.DecodeFile(path, config); err != nil {
		fmt.Println(err)
		panic(err)
	}
	return config
}
