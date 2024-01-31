package main

import (
	"log"
	v1 "movie-app/rating/internal/v1"
)

func main() {
	log.Println("Starting rating service")
	initApp("gin", "3011")
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
