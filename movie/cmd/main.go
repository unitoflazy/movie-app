package main

import (
	"log"
	v1 "movie-app/movie/internal/v1"
)

func main() {
	log.Println("Starting movie service")
	initApp("gin", "8000", "http://localhost:8010", "http://localhost:8020")
}

func initApp(appType string, port string, metadataAdr string, ratingAdr string) {
	var err error

	if appType == "gin" {
		err = v1.NewGinApp(port, metadataAdr, ratingAdr).Run()
	}

	if err != nil {
		log.Fatalln("failed to start "+appType+" app: ", err.Error())
	}
}
