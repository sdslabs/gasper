package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sdslabs/SDS/docker"
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

	// Content of the config file at container level
	content, err := ioutil.ReadFile(
		gopath + "/src/github.com/sdslabs/SDS/configs/containerLevel/template." + app +
			".sdslabs.co.conf")
	if err != nil {
		return err
	}

	// Replace `template` by name of the application
	conf := strings.Replace(string(content), "template", name, -1)
	content = []byte(conf)

	reader, err := TarFile(content, name+"."+app+".sdslabs.co.conf", 0644)
	if err != nil {
		return err
	}

	// Add the config file to the corresponding container
	err = docker.AddFileToContainer(containerID, "/etc/nginx/conf.d", reader)
	if err != nil {
		return err
	}

	return nil
}
