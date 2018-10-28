package utils

import (
	"archive/tar"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// TarFile function takes a source file and converts it into a tar archive
func TarFile(content []byte, filename string, mode int64) (io.Reader, error) {
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

// TarDir function takes a source directory and converts it into a tar archive
func TarDir(source string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	ok := filepath.Walk(source, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("1. %s", err)
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return fmt.Errorf("2. %s", err)
		}

		header.Name = strings.TrimPrefix(strings.Replace(file, source, "", -1), string(filepath.Separator))

		err = tw.WriteHeader(header)
		if err != nil {
			return fmt.Errorf("3. %s", err)
		}

		f, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("4. %s", err)
		}

		if fi.IsDir() {
			return nil
		}

		_, err = io.Copy(tw, f)
		if err != nil {
			return fmt.Errorf("5. %s", err)
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
