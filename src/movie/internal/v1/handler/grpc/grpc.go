package grpchandler

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movie-app/gen"
	"movie-app/metadata/pkg/model"
	"movie-app/movie/internal/v1/controller"
	"movie-app/movie/internal/v1/gateway"
)

type GRPCHandler struct {
	gen.UnimplementedMovieServiceServer
	ctrl *controller.Controller
}

func NewGRPCHandler(ctrl *controller.Controller) *GRPCHandler {
	return &GRPCHandler{ctrl: ctrl}
}

func (h *GRPCHandler) GetMovieDetails(ctx context.Context, req *gen.GetMovieDetailsRequest) (*gen.GetMovieDetailsResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	m, err := h.ctrl.GetMovieDetails(ctx, req.MovieId)
	if errors.Is(err, gateway.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetMovieDetailsResponse{
		MovieDetails: &gen.MovieDetails{
			Metadata: model.MetadataToProto(&m.Metadata),
			Rating:   *m.Rating,
		},
	}, nil
}
