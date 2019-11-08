package factory

import (
	"context"
	"errors"

	"github.com/sdslabs/gasper/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type credentials struct {
	Secret string
}

const secretKeyHolder = "secret"

func (c *credentials) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		secretKeyHolder: c.Secret,
	}, nil
}

func (c *credentials) RequireTransportSecurity() bool {
	return false
}

func authorize(ctx context.Context) error {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if len(md[secretKeyHolder]) == 0 {
			return errors.New("GRPC: Secret Key was not provided")
		}
		if md[secretKeyHolder][0] == configs.GasperConfig.Secret {
			return nil
		}
		return errors.New("GRPC: Invalid Secret Key")
	}
	return errors.New("GRPC: Failed to extract metadata")
}

func streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := authorize(stream.Context()); err != nil {
		return err
	}
	return handler(srv, stream)
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := authorize(ctx); err != nil {
		return nil, err
	}
	return handler(ctx, req)
}
