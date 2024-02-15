package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movie-app/gen"
	"movie-app/metadata/internal/v1/controller/metadata"
	"movie-app/metadata/internal/v1/repository"
	"movie-app/metadata/pkg/model"
)

type GRPCHandler struct {
	gen.UnimplementedMetadataServiceServer
	ctrl *metadata.Controller
}

func NewGRPCHandler(ctrl *metadata.Controller) *GRPCHandler {
	return &GRPCHandler{ctrl: ctrl}
}

func (h *GRPCHandler) GetMetadata(ctx context.Context, req *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	m, err := h.ctrl.Get(ctx, req.MovieId)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &gen.GetMetadataResponse{
		Metadata: model.MetadataToProto(m),
	}, nil
}
