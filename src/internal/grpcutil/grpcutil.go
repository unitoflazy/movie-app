package grpcutil

import (
	"context"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"movie-app/pkg/discovery"
)

func ServiceConnection(ctx context.Context, serviceName string, registry discovery.Registry, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	addrs, err := registry.ServiceAddresses(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	return grpc.Dial(addrs[rand.Intn(len(addrs))], opts...)
}
