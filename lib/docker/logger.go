package docker

import (
	"encoding/binary"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// ReadLogs returns the logs from a docker container
func ReadLogs(ctx context.Context, cli *client.Client, containerID, tail string) ([]string, error) {
	reader, err := cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: true,
		Tail:       tail,
	})

	defer reader.Close()

	if err != nil {
		return nil, err
	}

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
