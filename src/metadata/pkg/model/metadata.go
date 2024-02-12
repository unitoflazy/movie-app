package model

// Metadata is a struct that represents the rating of a api
type Metadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Director    string `json:"director"`
}
