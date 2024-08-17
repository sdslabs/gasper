package docker

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/sdslabs/gasper/lib/utils"
	"golang.org/x/net/context"
)

// Check for available images and pull if not present
func CheckAndPullImages(imageList ...string) {
	availableImages, err := ListImages()
	if err != nil {
		utils.LogError("Main-Helper-1", err)
		os.Exit(1)
	}
	for _, image := range imageList {
		imageWithoutRepoName := strings.Replace(image, "docker.io/", "", -1)
		if utils.Contains(availableImages, image) || utils.Contains(availableImages, imageWithoutRepoName) {
			continue
		}
		utils.LogInfo("Main-Helper-2", "Image %s not present locally, pulling from DockerHUB", image)
		if err = DirectPull(image); err != nil {
			utils.LogError("Main-Helper-3", err)
		}
	}
}

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
