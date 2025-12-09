package factory

import (
	"context"

	pb "github.com/sdslabs/gasper/lib/factory/protos/application"
	"github.com/sdslabs/gasper/types"
	"google.golang.org/grpc"
)

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
