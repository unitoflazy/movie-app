package grpchandler

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movie-app/gen"
	"movie-app/rating/internal/v1/controller/rating"
	"movie-app/rating/internal/v1/repository"
	"movie-app/rating/pkg/model"
)

type GRPCHandler struct {
	gen.UnimplementedRatingServiceServer
	ctrl *rating.Controller
}

func NewGRPCHandler(ctrl *rating.Controller) *GRPCHandler {
	return &GRPCHandler{ctrl: ctrl}
}

func (h *GRPCHandler) GetAggregatedRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "record id and record type are required")
	}

	r, err := h.ctrl.GetAggregatedRating(ctx, model.RecordType(req.RecordType), model.RecordID(req.RecordId))
	if errors.Is(err, repository.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetAggregatedRatingResponse{
		RatingValue: r,
	}, nil
}

func (h *GRPCHandler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "record id record type are required")
	}

	if err := h.ctrl.PutRating(ctx, model.RecordType(req.RecordType), model.RecordID(req.RecordId), &model.Rating{
		Value:  model.RatingValue(req.RatingValue),
		UserID: model.UserID(req.UserId),
	}); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.PutRatingResponse{}, nil
}
