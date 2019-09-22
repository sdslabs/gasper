package docker

import (
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sdslabs/SWS/lib/utils"
	"golang.org/x/net/context"
)

// ListImages function returns a list of docker images present in the system
func ListImages() []string {
	cli, err := client.NewEnvClient()
	if err != nil {
		utils.LogError(err)
		panic(err)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		utils.LogError(err)
		panic(err)
	}

	list := make([]string, 1)

	for _, image := range images {
		if len(image.RepoTags) > 0 {
			list = append(list, image.RepoTags[0])
		}
	}
	return list
}

// Pull function pulls an image from DockerHUB
func Pull(image string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		utils.LogError(err)
		panic(err)
	}

	out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		utils.LogError(err)
		panic(err)
	}

	defer out.Close()

	io.Copy(os.Stdout, out)
}
