package memory

import (
	"context"
	"movie-app/rating/internal/v1/repository"
	"movie-app/rating/pkg/model"
)

type RatingMap map[model.RecordType]map[model.RecordID][]model.Rating

type Repository struct {
	data RatingMap
}

func New() *Repository {
	return &Repository{
		data: make(RatingMap),
	}
}

func (r *Repository) Get(ctx context.Context, recordType model.RecordType, recordID model.RecordID) ([]model.Rating, error) {
	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}

	if _, ok := r.data[recordType][recordID]; !ok {
		return nil, repository.ErrNotFound
	}

	return r.data[recordType][recordID], nil
}

func (r *Repository) Put(ctx context.Context, recordType model.RecordType, recordID model.RecordID, rating *model.Rating) error {
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = make(map[model.RecordID][]model.Rating)
	}

	if _, ok := r.data[recordType][recordID]; !ok {
		r.data[recordType][recordID] = make([]model.Rating, 0)
	}

	r.data[recordType][recordID] = append(r.data[recordType][recordID], *rating)

	return nil
}
