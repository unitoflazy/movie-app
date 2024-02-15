package model

type RateMovieRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Value  int    `json:"rating_value" binding:"required"`
}
