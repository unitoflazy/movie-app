package model

type RateMovieRequest struct {
	RecordID   RecordID   `json:"recordId" binding:"required"`
	RecordType RecordType `json:"recordType" binding:"required"`
	UserID     string     `json:"userId" binding:"required"`
	Value      int        `json:"value" binding:"required"`
}
