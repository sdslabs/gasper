package factory

import (
	"context"
	"encoding/json"

	pb "github.com/sdslabs/gasper/lib/factory/protos/application"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
	"google.golang.org/grpc"
)

// GracefulUp starts all stopped application containers in the given worker node, and creates apps whose containers are not present
func GracefulUp(instanceURL string, apps []types.ApplicationConfig) (bool, error) {
	conn, err := grpc.Dial(
		instanceURL,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(authCredentials),
	)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	client := pb.NewApplicationFactoryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var appNames []string

	for _, app := range apps {
		appNames = append(appNames, app.Name)
	}

	res, err := client.FetchRunningContainers(ctx, &pb.NodeInventory{
		InstanceURL:   instanceURL,
		ContainerList: appNames,
	})
	if err != nil {
		return false, err
	}

	for _, app := range apps {
		if !utils.Contains(res.ContainerList, app.Name) {
			utils.LogInfo("APP", "Creating App %s on node %s", app.Name, instanceURL)
			data, err := json.Marshal(app)
			if err != nil {
				return false, err
			}
			_, err = client.Create(ctx, &pb.RequestBody{
				Language: app.Language,
				Owner:    app.Owner,
				Data:     data,
			})
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

func GracefulDown(instanceURL string, apps []types.ApplicationConfig) (bool, error) {
	conn, err := grpc.Dial(
		instanceURL,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(authCredentials),
	)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	client := pb.NewApplicationFactoryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, app := range apps {
		_, err := client.RemoveContainer(ctx, &pb.NameHolder{Name: app.Name})
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
