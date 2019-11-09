package factory

import (
	"context"

	pb "github.com/sdslabs/gasper/lib/factory/protos/database"
	"google.golang.org/grpc"
)

// CreateDatabase is a remote procedure call for creating a database in a worker node
func CreateDatabase(language, owner, instanceURL string, data []byte) ([]byte, error) {
	conn, err := grpc.Dial(
		instanceURL,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(authCredentials),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewDatabaseFactoryClient(conn)

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

// DeleteDatabase is a remote procedure call for deleting a database in a worker node
func DeleteDatabase(name, instanceURL string) (*pb.GenericResponse, error) {
	conn, err := grpc.Dial(
		instanceURL,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(authCredentials),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewDatabaseFactoryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := client.Delete(ctx, &pb.NameHolder{Name: name})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FetchDatabaseServerLogs is a remote procedure call for fetching logs of a database server in a worker node
func FetchDatabaseServerLogs(language, tail, instanceURL string) (*pb.LogResponse, error) {
	conn, err := grpc.Dial(
		instanceURL,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(authCredentials),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewDatabaseFactoryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := client.FetchLogs(ctx, &pb.LogRequest{
		Language: language,
		Tail:     tail,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ReloadDatabaseServer is a remote procedure call for restarting a database server in a worker node
func ReloadDatabaseServer(language, instanceURL string) (*pb.GenericResponse, error) {
	conn, err := grpc.Dial(
		instanceURL,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(authCredentials),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewDatabaseFactoryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := client.Reload(ctx, &pb.LanguageHolder{Language: language})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// NewDatabaseFactory returns a new GRPC server for creating databases
func NewDatabaseFactory(bindings pb.DatabaseFactoryServer) *grpc.Server {
	srv := grpc.NewServer(
		grpc.StreamInterceptor(streamInterceptor),
		grpc.UnaryInterceptor(unaryInterceptor),
	)
	pb.RegisterDatabaseFactoryServer(srv, bindings)
	return srv
}
