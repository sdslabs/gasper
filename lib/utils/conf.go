package utils

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/sdslabs/SWS/lib/docker"
)

// ServiceConfFunc is a type of function that takes the service app name and
// returns the conf file for the server
type ServiceConfFunc func(name string) string

// WriteServiceConfFile takes the name of service app and config func and
// writes the conf file of server in the required container
func WriteServiceConfFile(ctx context.Context, cli *client.Client, containerID, name string, createServiceConf ServiceConfFunc) error {
	conf := []byte(createServiceConf(name))
	filename := fmt.Sprintf("%s.sws.conf", name)
	reader, err := NewTarArchiveFromContent(conf, filename, 0644)
	if err != nil {
		return err
	}
	err = docker.CopyToContainer(ctx, cli, containerID, "/etc/nginx/conf.d/", reader)
	if err != nil {
		return err
	}
	return nil
}
