package configs

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

func getConfiguration() *GasperCfg {
	flag.Parse()
	config := &GasperCfg{}
	if _, err := toml.DecodeFile(*configFile, config); err != nil {
		fmt.Println("\x1b[35m[\x1b[0m\x1b[31mERROR\x1b[0m\x1b[35m]\x1b[0m \x1b[91m>>>\x1b[0m \x1b[32m", err.Error(), "\x1b[0m")
		os.Exit(1)
	}
	return config
}
