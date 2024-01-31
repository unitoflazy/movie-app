package main

import (
	"log"
	v1 "movie-app/metadata/internal/v1"
)

func main() {
	log.Println("Starting movie metadata service")
	initApp("gin", "3012")
}

func initApp(appType, port string) {
	var err error = nil

	if appType == "gin" {
		err = v1.NewGinApp(port).Run()
	}

	if err != nil {
		log.Fatalln("failed to start "+appType+" app: ", err.Error())
	}
}
