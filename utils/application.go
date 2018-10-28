package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sdslabs/SDS/docker"
	git "gopkg.in/src-d/go-git.v4"
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
		"%s/src/github.com/sdslabs/SDS/configs/containerLevel/template.%s.sdslabs.co.conf",
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

	stream, err := TarFile(content, targetFile, 0644)

	// Add the config file to the corresponding container
	err = docker.AddToContainer(containerID, "/etc/nginx/conf.d/", stream)
	if err != nil {
		return err
	}

	return nil
}

// GitCloneApp takes the github url and clones it into the given container
func GitCloneApp(name, url, containerID string) error {
	_, err := os.Stat("/tmp/SDS")
	if os.IsNotExist(err) {
		os.Mkdir("/tmp/SDS", 0755)
	}

	// Plain clone the repo in `/tmp/SDS/<name>`
	dirname := fmt.Sprintf("/tmp/SDS/%s", name)

	_, err = git.PlainClone(dirname, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return err
	}

	// Tar folder and copy to container
	stream, err := TarDir(dirname)

	err = docker.AddToContainer(containerID, fmt.Sprintf("/SDS/%s", name), stream)
	if err != nil {
		return err
	}

	return nil
}
