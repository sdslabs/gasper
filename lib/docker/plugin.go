package docker

import (
	"io"

	dockerTypes "github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

//ListPlugins lists all plugins
func ListPlugins() ([]string, error) {
	ctx := context.Background()
	plugins, err := cli.PluginList(ctx)
	if err != nil {
		return nil, err
	}

	list := make([]string, 0)

	for _, plugin := range plugins {
		if len(plugin.Name) > 0 {
			list = append(list, plugin.Name)
		}
	}
	return list, nil
}

//InstallPlugin installs required plugin
func InstallPlugin(name string, options dockerTypes.PluginInstallOptions) (io.ReadCloser, error) {
	ctx := context.Background()
	rc, err := cli.PluginInstall(ctx, name, options)
	return rc, err
}

//EnablePlugin enables required plugin
func EnablePlugin(name string) error {
	ctx := context.Background()
	err := cli.PluginEnable(ctx, name, dockerTypes.PluginEnableOptions{Timeout: 30})
	return err
}

//IsPluginEnabled checks whether plugin is enabled
func IsPluginEnabled(name string) (bool, error) {
	ctx := context.Background()
	plugin, _, err := cli.PluginInspectWithRaw(ctx, name)
	if err != nil {
		return false, err
	}
	return plugin.Enabled, err
}
