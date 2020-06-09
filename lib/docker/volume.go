package docker

import (
	dockerfilter "github.com/docker/docker/api/types/filters"
	volumetypes "github.com/docker/docker/api/types/volume"
	"golang.org/x/net/context"
)

// ListVolumes lists all the volumes
func ListVolumes() ([]string, error) {
	ctx := context.Background()
	volumestruct, err := cli.VolumeList(ctx, dockerfilter.Args{})
	if err != nil {
		return nil, err
	}

	volumes := volumestruct.Volumes

	list := make([]string, 0)

	for _, volume := range volumes {
		if len(volume.Name) > 0 {
			list = append(list, volume.Name)
		}
	}
	return list, nil
}

// CreateVolume creates a volume with given name and driver
func CreateVolume(name string, driver string) (string, error) {
	ctx := context.Background()
	volume, err := cli.VolumeCreate(ctx, volumetypes.VolumesCreateBody{Driver: driver, Labels: map[string]string{}, Name: name, DriverOpts: map[string]string{"ReplicationGoal": ""}})
	return volume.Name, err
}
