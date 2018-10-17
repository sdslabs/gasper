package utils

import (
	"archive/tar"
	"bytes"
	"io"
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

	return tar.NewReader(&buf), nil
}
