package factory

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/sdslabs/gasper/lib/factory/protos/application"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
	"google.golang.org/grpc"
)

// GracefulUp starts all stopped application containers in the given worker node
func GracefulUp(instanceURL string, appsOnNode []types.ApplicationConfig) (bool, error) {
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

	res, err := client.StartStoppedAppContainers(ctx, &pb.ContainerRequestBody{
		InstanceURL: instanceURL,
	})
	if err != nil {
		return false, err
	}

	for _, app := range appsOnNode {
		if !utils.Contains(res.Data, app.Name) {
			fmt.Println("Creating app ", app.Name)
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
