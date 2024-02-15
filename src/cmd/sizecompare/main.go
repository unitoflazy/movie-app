package main

import (
	"encoding/xml"
	"fmt"
	json "github.com/bytedance/sonic"
	"github.com/golang/protobuf/proto"
	"movie-app/gen"
	"movie-app/metadata/pkg/model"
)

var metadata = &model.Metadata{
	ID:          "123",
	Title:       "The Dark Knight",
	Description: "Batman fights crime in Gotham City",
	Director:    "Christopher Nolan",
}

var genMetadata = &gen.Metadata{
	Id:          "123",
	Title:       "The Dark Knight",
	Description: "Batman fights crime in Gotham City",
	Director:    "Christopher Nolan",
}

func main() {
	jsonBytes, err := serializeToJSON(metadata)
	if err != nil {
		panic(err)
	}

	xmlBytes, err := serializeToXML(metadata)
	if err != nil {
		panic(err)
	}

	protoBytes, err := serializeToProto(genMetadata)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JSON size: %d\n", len(jsonBytes))
	fmt.Printf("XML size: %d\n", len(xmlBytes))
	fmt.Printf("Proto size: %d\n", len(protoBytes))
}

func serializeToJSON(metadata *model.Metadata) ([]byte, error) {
	return json.Marshal(metadata)
}

func serializeToXML(metadata *model.Metadata) ([]byte, error) {
	return xml.Marshal(metadata)
}

func serializeToProto(metadata *gen.Metadata) ([]byte, error) {
	return proto.Marshal(metadata)
}
