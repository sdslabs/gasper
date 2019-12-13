package docker

import (
	"encoding/binary"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ReadLogs returns the logs from a docker container
func ReadLogs(containerID, tail string) ([]string, error) {
	ctx := context.Background()
	reader, err := cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: true,
		Tail:       tail,
	})

	if err != nil {
		return nil, err
	}

	defer reader.Close()

	logs := []string{}
	hdr := make([]byte, 8)

	for {
		_, err = reader.Read(hdr)

		if err != nil {
			return logs, err
		}

		count := binary.BigEndian.Uint32(hdr[4:])
		dat := make([]byte, count)
		_, err = reader.Read(dat)
		logs = append(logs, string(dat))

		if err != nil {
			return logs, err
		}
	}
}
