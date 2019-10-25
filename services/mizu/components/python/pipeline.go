package python

import (
	"fmt"
	"strings"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/api"
	"github.com/sdslabs/gasper/lib/commons"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/types"
)

const (
	pythonVersionTag = "python_version"
	python3Tag       = "3"
	python2Tag       = "2"
)

func bootstrap(requirementsPath, index string, args []string, env *types.ApplicationEnv) (string, types.ResponseError) {
	arguments := strings.Join(args, " ")
	var serveCmd string
	if requirementsPath == "" {
		serveCmd = fmt.Sprintf(`python %s %s &> /proc/1/fd/1`, index, arguments)
	} else {
		serveCmd = fmt.Sprintf(`pip install -r %s &> /proc/1/fd/1; python %s %s &> /proc/1/fd/1`, requirementsPath, index, arguments)
	}
	cmd := []string{"sh", "-c", serveCmd}
	execID, err := docker.ExecDetachedProcess(env.ContainerID, cmd)
	if err != nil {
		return execID, types.NewResErr(500, "failed to start the server", err)
	}
	return execID, nil
}

// Pipeline is the application creation pipeline
func Pipeline(data types.M) types.ResponseError {
	var image string
	if data[pythonVersionTag].(string) == python3Tag {
		image = configs.ImageConfig.Python3
	} else if data[pythonVersionTag].(string) == python2Tag {
		image = configs.ImageConfig.Python2
	}

	appConf := &types.ApplicationConfig{
		DockerImage:  image,
		ConfFunction: configs.CreatePythonContainerConfig,
	}

	appEnv, resErr := api.SetupApplication(appConf, data)
	if resErr != nil {
		return resErr
	}

	context := data["context"].(map[string]interface{})

	if context["rcFile"].(bool) {
		return nil
	}

	// Path of `requirements.txt` or any-other file containing requirements
	var requirements = ""
	if data["requirements"] != nil {
		requirements = data["requirements"].(string)
	}

	if data["django"] == nil {
		data["django"] = false
	}

	if data["django"].(bool) {
		_, resErr = bootstrap(requirements, "manage.py", []string{"runserver"}, appEnv)
	} else {
		var args []interface{}
		if context["args"] != nil {
			args = context["args"].([]interface{})
		}
		var arguments []string
		for _, arg := range args {
			arguments = append(arguments, arg.(string))
		}
		_, resErr = bootstrap(requirements, context["index"].(string), arguments, appEnv)
	}
	if resErr != nil {
		go commons.AppFullCleanup(data["name"].(string))
		return resErr
	}

	return nil
}
