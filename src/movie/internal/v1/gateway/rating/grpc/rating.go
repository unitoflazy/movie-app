package grpc

import (
	"context"
	"google.golang.org/grpc"
	"movie-app/gen"
	"movie-app/internal/grpcutil"
	"movie-app/pkg/discovery"
	ratingmodel "movie-app/rating/pkg/model"
)

type Gateway struct {
	registry discovery.Registry
	opts     []grpc.DialOption
}

func New(registry *discovery.Registry, opts ...grpc.DialOption) *Gateway {
	return &Gateway{registry: *registry, opts: opts}
}

func (g *Gateway) GetAggregatedRating(ctx context.Context, ratingType ratingmodel.RecordType, id ratingmodel.RecordID) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry, g.opts...)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := gen.NewRatingServiceClient(conn)
	resp, err := client.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{
		RecordId:   string(id),
		RecordType: string(ratingType),
	})
	if err != nil {
		return 0, err
	}

	return resp.RatingValue, nil
}

func (g *Gateway) PutRating(ctx context.Context, ratingType ratingmodel.RecordType, id ratingmodel.RecordID, rating *ratingmodel.Rating) error {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry, g.opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := gen.NewRatingServiceClient(conn)
	_, err = client.PutRating(ctx, &gen.PutRatingRequest{
		RecordId:    string(id),
		RecordType:  string(ratingType),
		RatingValue: int32(rating.Value),
	})

	return err
}
