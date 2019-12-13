package factory

import (
	"context"

	pb "github.com/sdslabs/gasper/lib/factory/protos/application"
	"google.golang.org/grpc"
)

// CreateApplication is a remote procedure call for creating an application in a worker node
func CreateApplication(language, owner, instanceURL string, data []byte) ([]byte, error) {
	conn, err := grpc.Dial(
		instanceURL,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(authCredentials),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewApplicationFactoryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := client.Create(ctx, &pb.RequestBody{
		Language: language,
		Owner:    owner,
		Data:     data,
	})
	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}

// RebuildApplication is a remote procedure call for rebuilding an application in a worker node
func RebuildApplication(name, instanceURL string) ([]byte, error) {
	conn, err := grpc.Dial(
		instanceURL,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(authCredentials),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewApplicationFactoryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := client.Rebuild(ctx, &pb.NameHolder{Name: name})
	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}

// DeleteApplication is a remote procedure call for deleting an application in a worker node
func DeleteApplication(name, instanceURL string) (*pb.DeletionResponse, error) {
	conn, err := grpc.Dial(
		instanceURL,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(authCredentials),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewApplicationFactoryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := client.Delete(ctx, &pb.NameHolder{Name: name})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FetchApplicationLogs is a remote procedure call for fetching logs of an application in a worker node
func FetchApplicationLogs(name, tail, instanceURL string) (*pb.LogResponse, error) {
	conn, err := grpc.Dial(
		instanceURL,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(authCredentials),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewApplicationFactoryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := client.FetchLogs(ctx, &pb.LogRequest{
		Name: name,
		Tail: tail,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// NewApplicationFactory returns a new GRPC server for creating applications
func NewApplicationFactory(bindings pb.ApplicationFactoryServer) *grpc.Server {
	srv := grpc.NewServer(
		grpc.StreamInterceptor(streamInterceptor),
		grpc.UnaryInterceptor(unaryInterceptor),
	)
	pb.RegisterApplicationFactoryServer(srv, bindings)
	return srv
}
