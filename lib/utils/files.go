package utils

import (
	"archive/tar"
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// NewTarArchiveFromContent returns an io.Reader type tar archive from content in []byte
func NewTarArchiveFromContent(content []byte, filename string, mode int64) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	header := &tar.Header{
		Name: filename,
		Mode: mode,
		Size: int64(len(content)),
	}
	err := tw.WriteHeader(header)
	if err != nil {
		return nil, err
	}

	_, err = tw.Write(content)
	if err != nil {
		return nil, err
	}
	err = tw.Close()
	if err != nil {
		return nil, err
	}

	return bufio.NewReader(&buf), nil
}

// NewTarArchiveFromPath returns an io.Reader type tar archive from content in given path
func NewTarArchiveFromPath(path string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	ok := filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(strings.Replace(file, path, "", -1), string(filepath.Separator))
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		_, err = io.Copy(tw, f)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
		return nil
	})

	if ok != nil {
		return nil, ok
	}
	ok = tw.Close()
	if ok != nil {
		return nil, ok
	}
	return bufio.NewReader(&buf), nil
}
