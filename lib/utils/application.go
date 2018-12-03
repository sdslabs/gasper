package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sdslabs/SWS/lib/docker"
)

// ReadAndWriteConfig takes the template config and writes to the corresponding container
// Takes in three arguments - name and app (`static`, `php`, `node`, 'python' and `go`)
// and corresponding containerID
func ReadAndWriteConfig(name string, app string, containerID string) error {
	// Set gopath to absolute gopath in the environment
	gopath, ok := os.LookupEnv("GOPATH")
	if !ok {
		return errors.New("Environment variable `GOPATH` does not exist")
	}

	fileName := fmt.Sprintf(
		"%s/src/github.com/sdslabs/SWS/configs/containerLevel/template.%s.sdslabs.co.conf",
		gopath, app)

	// Content of the config file at container level
	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		return err
	}

	// Replace `template` by name of the application
	conf := strings.Replace(string(file), "template", name, -1)

	content := []byte(conf)
	targetFile := fmt.Sprintf("%s.%s.sdslabs.co.conf", name, app)

	stream, err := tarFile(content, targetFile, 644)

	// Add the config file to the corresponding container
	err = docker.AddFileToContainer(containerID, "/etc/nginx/conf.d/", stream)
	if err != nil {
		return err
	}

	return nil
}
