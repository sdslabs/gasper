package main

import (
	"fmt"
	"os/exec"
)

var releaseMap = map[string][]string{
	"darwin":  []string{"amd64"},
	"linux":   []string{"386", "amd64", "arm", "arm64", "ppc64le"},
	"freebsd": []string{"386", "amd64"},
	"openbsd": []string{"amd64"},
	"windows": []string{"386", "amd64"},
}

func build(platform, architecture string, errChan chan error) {
	_, err := exec.Command("sh", "-c",
		fmt.Sprintf("GOOS=%s GOARCH=%s go build -o bin/gasper_%s_%s", platform, architecture, platform, architecture)).Output()
	errChan <- err
}

func main() {
	syncGroup := make([]chan error, 0)
	for platform, architectureList := range releaseMap {
		for _, architecture := range architectureList {
			errChan := make(chan error)
			syncGroup = append(syncGroup, errChan)
			go build(platform, architecture, errChan)
		}
	}
	for _, errChan := range syncGroup {
		if err := <-errChan; err != nil {
			fmt.Println(err)
		}
	}
}
