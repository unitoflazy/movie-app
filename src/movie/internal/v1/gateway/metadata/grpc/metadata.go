package grpc

import (
	"context"
	"google.golang.org/grpc"
	"movie-app/gen"
	"movie-app/internal/grpcutil"
	"movie-app/metadata/pkg/model"
	"movie-app/pkg/discovery"
)

type Gateway struct {
	registry discovery.Registry
	opts     []grpc.DialOption
}

func New(registry *discovery.Registry, opts ...grpc.DialOption) *Gateway {
	return &Gateway{registry: *registry, opts: opts}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry, g.opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})
	if err != nil {
		return nil, err
	}

	return model.ProtoToMetadata(resp.Metadata), nil
}
