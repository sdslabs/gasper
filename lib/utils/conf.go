package utils

import (
	"fmt"

	"github.com/sdslabs/SWS/lib/docker"
)

// ServiceConfFunc is a type of function that takes the service app name and
// returns the conf file for the server
type ServiceConfFunc func(name string) string

// WriteServiceConfFile takes the name of service app and config func and
// writes the conf file of server in the required container
func WriteServiceConfFile(containerID, name string, createServiceConf ServiceConfFunc) error {
	conf := []byte(createServiceConf(name))
	filename := fmt.Sprintf("%s.sws.conf", name)
	stream, err := TarFile(conf, filename, 0644)
	if err != nil {
		return err
	}
	err = docker.CopyToContainer(containerID, "/etc/nginx/conf.d/", stream)
	if err != nil {
		return err
	}
	return nil
}
