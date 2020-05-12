package docker

import (
	"io"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ListImages function returns a list of docker images present in the system
func ListImages() ([]string, error) {
	ctx := context.Background()
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}

	list := make([]string, 1)

	for _, image := range images {
		if len(image.RepoTags) > 0 {
			list = append(list, image.RepoTags[0])
		}
	}
	return list, nil
}

// Pull function pulls an image from DockerHUB
func Pull(image string) error {
	ctx := context.Background()
	out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()

	io.Copy(os.Stdout, out)
	return nil
}

// DirectPull function directly pulls an image from DockerHUB using os/exec
func DirectPull(image string) error {
	cmd := exec.Command("docker", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
