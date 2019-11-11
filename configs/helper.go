package configs

import (
	"flag"
	"fmt"

	"github.com/BurntSushi/toml"
)

func getConfiguration() *GasperCfg {
	flag.Parse()
	config := &GasperCfg{}
	if _, err := toml.DecodeFile(*configFile, config); err != nil {
		fmt.Println(err)
		panic(err)
	}
	return config
}
