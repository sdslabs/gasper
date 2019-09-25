package python

import (
	"fmt"
	"strings"

	"github.com/sdslabs/SWS/configs"
	"github.com/sdslabs/SWS/lib/api"
	"github.com/sdslabs/SWS/lib/commons"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
)

const (
	pythonVersionTag = "python_version"
	python3Tag       = "3"
	python2Tag       = "2"
)

func startServer(index string, args []string, env *types.ApplicationEnv) (string, types.ResponseError) {
	arguments := strings.Join(args, " ")
	serveCmd := fmt.Sprintf(`python %s %s &> /proc/1/fd/1`, index, arguments)
	cmd := []string{"bash", "-c", serveCmd}
	execID, err := docker.ExecDetachedProcess(env.Context, env.Client, env.ContainerID, cmd)
	if err != nil {
		return execID, types.NewResErr(500, "failed to start the server", err)
	}
	return execID, nil
}

func installRequirements(path string, env *types.ApplicationEnv) (string, types.ResponseError) {
	cmd := []string{"bash", "-c", fmt.Sprintf(`pip install -r %s &> /proc/1/fd/1`, path)}
	execID, err := docker.ExecDetachedProcess(env.Context, env.Client, env.ContainerID, cmd)
	if err != nil {
		return execID, types.NewResErr(500, "failed to install requirements", err)
	}
	return execID, nil
}

// Pipeline is the application creation pipeline
func Pipeline(data map[string]interface{}) types.ResponseError {
	var image string
	if data[pythonVersionTag].(string) == python3Tag {
		image = configs.ServiceConfig["python"].(map[string]interface{})["python3_image"].(string)
	} else if data[pythonVersionTag].(string) == python2Tag {
		image = configs.ServiceConfig["python"].(map[string]interface{})["python2_image"].(string)
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
	requirements := data["requirements"]
	if requirements != nil {
		_, resErr = installRequirements(requirements.(string), appEnv)
		if resErr != nil {
			go commons.FullCleanup(data["name"].(string), data["instanceType"].(string))
			return resErr
		}
	}

	if data["django"] == nil {
		data["django"] = false
	}

	if data["django"].(bool) {
		_, resErr = startServer("manage.py", []string{"runserver"}, appEnv)
	} else {
		var args []interface{}
		if context["args"] != nil {
			args = context["args"].([]interface{})
		}
		var arguments []string
		for _, arg := range args {
			arguments = append(arguments, arg.(string))
		}
		_, resErr = startServer(context["index"].(string), arguments, appEnv)
	}
	if resErr != nil {
		go commons.FullCleanup(data["name"].(string), data["instanceType"].(string))
		return resErr
	}

	return nil
}
