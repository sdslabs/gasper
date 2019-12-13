package factory

import (
	"time"

	"github.com/sdslabs/gasper/configs"
)

const timeout = 30 * time.Second

var authCredentials = &credentials{Secret: configs.GasperConfig.Secret}
