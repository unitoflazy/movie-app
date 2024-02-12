package model

type RateMovieRequest struct {
	UserID string `json:"userId" binding:"required"`
	Value  int    `json:"value" binding:"required"`
}
