package model

import "movie-app/gen"

func MetadataToProto(metadata *Metadata) *gen.Metadata {
	return &gen.Metadata{
		Id:          metadata.ID,
		Title:       metadata.Title,
		Description: metadata.Description,
		Director:    metadata.Director,
	}
}

func ProtoToMetadata(metadata *gen.Metadata) *Metadata {
	return &Metadata{
		ID:          metadata.Id,
		Title:       metadata.Title,
		Description: metadata.Description,
		Director:    metadata.Director,
	}
}
