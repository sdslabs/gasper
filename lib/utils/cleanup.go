package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdslabs/SWS/lib/docker"
)

// StorageCleanup removes the application's local storage directory
func StorageCleanup(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println(err)
	}
}

// ContainerCleanup removes the application's container
func ContainerCleanup(appName string) {
	docker.DeleteContainer(appName)
}

// FullCleanup cleans the specified application's container and local storage
func FullCleanup(appName string) {
	var (
		path, _ = os.Getwd()
		appDir  = filepath.Join(path, fmt.Sprintf("storage/%s", appName))
	)
	StorageCleanup(appDir)
	ContainerCleanup(appName)
}
