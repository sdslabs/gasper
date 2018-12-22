package php

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
	"golang.org/x/net/context"
)

func installPackages(composer bool) *types.ResponseError {
	if composer {
		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			return types.NewResponseError(500, "", err)
		}
		cmd := []string{"bash", "-c", "curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer"}
		execId, err := docker.ExecDetachedProcess(ctx, cli, "c859d83d6ac1", cmd)
		fmt.Println(execId)
		if err != nil {
			return types.NewResponseError(500, "failed to install composer in the container", err)
		}
	}
	return nil
}
