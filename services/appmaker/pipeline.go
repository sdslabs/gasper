package appmaker

import (
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/api"
	"github.com/sdslabs/gasper/types"
)

// applicationHandler is a struct for managing the creation of a new application of a specific language
type applicationHandler struct {
	image         string
	confGenerator func(string, string) string
}

// create handles the creation of a new application
func (handler *applicationHandler) create(app *types.ApplicationConfig) types.ResponseError {
	app.SetDockerImage(handler.image)
	app.SetConfGenerator(handler.confGenerator)
	return api.SetupApplication(app)
}

var pipeline = map[string]*applicationHandler{
	"nodejs": {
		image: configs.ImageConfig.Nodejs,
	},
	"python2": {
		image: configs.ImageConfig.Python2,
	},
	"python3": {
		image: configs.ImageConfig.Python3,
	},
	"golang": {
		image: configs.ImageConfig.Golang,
	},
	"ruby": {
		image: configs.ImageConfig.Ruby,
	},
	"rust": {
		image: configs.ImageConfig.Rust,
	},
	"php": {
		image:         configs.ImageConfig.Php,
		confGenerator: configs.CreatePHPContainerConfig,
	},
	"static": {
		image:         configs.ImageConfig.Static,
		confGenerator: configs.CreateStaticContainerConfig,
	},
}
