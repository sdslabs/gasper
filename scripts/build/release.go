package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

const releaseDir = "releases"

var releaseMap = map[string][]string{
	"darwin":  []string{"amd64"},
	"linux":   []string{"386", "amd64", "arm", "arm64", "ppc64le"},
	"freebsd": []string{"386", "amd64"},
	"openbsd": []string{"amd64"},
	"windows": []string{"386", "amd64"},
}

func md5sum(archiveName string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("md5 -r %s >> %s", archiveName, "checksums.txt"))
	cmd.Dir = releaseDir
	return cmd.Run()
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func build(platform, architecture, version string, errChan chan error) {
	fullName := fmt.Sprintf("gasper_%s_%s_%s", version, platform, architecture)
	dir := fmt.Sprintf("%s/%s", releaseDir, fullName)
	binaryName := "gasper"
	if platform == "windows" {
		binaryName += ".exe"
	}
	binaryPath := fmt.Sprintf("%s/%s", dir, binaryName)
	configPath := fmt.Sprintf("%s/%s", dir, "config.toml")
	if err := os.MkdirAll(dir, 0755); err != nil {
		errChan <- err
		return
	}
	if err := copyFile("config.sample.toml", configPath); err != nil {
		errChan <- err
		return
	}
	_, err := exec.Command("sh", "-c",
		fmt.Sprintf("GOOS=%s GOARCH=%s go build -o %s", platform, architecture, binaryPath)).Output()
	if err != nil {
		errChan <- err
		return
	}
	var archiveName string
	var cmd *exec.Cmd
	if platform == "windows" {
		archiveName = fullName + ".zip"
		cmd = exec.Command("zip", "-r", archiveName, fullName)
	} else {
		archiveName = fullName + ".tar.gz"
		cmd = exec.Command("tar", "-zcvf", archiveName, fullName)
	}
	cmd.Dir = releaseDir
	if err = cmd.Run(); err != nil {
		errChan <- err
		return
	}
	if err := md5sum(archiveName); err != nil {
		errChan <- err
		return
	}
	errChan <- os.RemoveAll(dir)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the VERSION argument")
		fmt.Println("Example:- make release VERSION=v1.0")
		return
	}
	version := os.Args[1]
	os.MkdirAll(releaseDir, 0755)
	syncGroup := make([]chan error, 0)
	for platform, architectureList := range releaseMap {
		for _, architecture := range architectureList {
			errChan := make(chan error)
			syncGroup = append(syncGroup, errChan)
			go build(platform, architecture, version, errChan)
		}
	}
	for _, errChan := range syncGroup {
		if err := <-errChan; err != nil {
			fmt.Println(err)
		}
	}
}
