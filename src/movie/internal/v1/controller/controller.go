package controller

import (
	"context"
	metadatamodel "movie-app/metadata/pkg/model"
	"movie-app/movie/pkg/model"
	ratingmodel "movie-app/rating/pkg/model"
)

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, ratingType ratingmodel.RecordType, id ratingmodel.RecordID) (float64, error)
	PutRating(ctx context.Context, ratingType ratingmodel.RecordType, id ratingmodel.RecordID, rating *ratingmodel.Rating) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadatamodel.Metadata, error)
}

type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{
		ratingGateway:   ratingGateway,
		metadataGateway: metadataGateway,
	}
}

func (c *Controller) GetMovieDetails(ctx context.Context, id string) (*model.MovieDetails, error) {
	md, err := c.metadataGateway.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	details := &model.MovieDetails{
		Metadata: *md,
	}

	rating, err := c.ratingGateway.GetAggregatedRating(ctx, ratingmodel.RecordTypeMovie, ratingmodel.RecordID(id))
	if err != nil {
		return nil, err
	} else {
		details.Rating = &rating
	}

	return details, nil
}
