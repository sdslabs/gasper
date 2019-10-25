package docker

import (
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/sdslabs/gasper/lib/utils"
	"golang.org/x/net/context"
)

// ListImages function returns a list of docker images present in the system
func ListImages() []string {
	ctx := context.Background()
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
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
	out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		utils.LogError(err)
		panic(err)
	}

	defer out.Close()

	io.Copy(os.Stdout, out)
}
