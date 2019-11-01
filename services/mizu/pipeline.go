package mizu

import (
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/api"
	"github.com/sdslabs/gasper/types"
)

func buildPipeline(dockerImage string, confGenerator func(string, string) string) func(*types.ApplicationConfig) types.ResponseError {
	return func(app *types.ApplicationConfig) types.ResponseError {
		app.SetDockerImage(dockerImage)
		app.SetConfGenerator(confGenerator)
		return api.SetupApplication(app)
	}
}

var pipeline = map[string]func(*types.ApplicationConfig) types.ResponseError{
	"nodejs":  buildPipeline(configs.ImageConfig.Nodejs, nil),
	"php":     buildPipeline(configs.ImageConfig.Php, configs.CreatePHPContainerConfig),
	"python2": buildPipeline(configs.ImageConfig.Python2, nil),
	"python3": buildPipeline(configs.ImageConfig.Python3, nil),
	"golang":  buildPipeline(configs.ImageConfig.Golang, nil),
	"static":  buildPipeline(configs.ImageConfig.Static, configs.CreateStaticContainerConfig),
}
