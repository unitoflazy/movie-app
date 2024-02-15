package model

// RecordID defines a record id. Together with RecordType
// identifies unique records across all types
type RecordID string

// RecordType defines a record type. Together with RecordID
// identifies unique records across all types
type RecordType string

const (
	RecordTypeMovie = RecordType("api")
)

type UserID string

// RatingValue defines a value of a rating record.
type RatingValue int

type Rating struct {
	RecordID   RecordID    `json:"record_id"`
	RecordType RecordType  `json:"record_type"`
	UserID     UserID      `json:"user_id"`
	Value      RatingValue `json:"value"`
}
