package main

import (
	"log"
	"movie-app/metadata/internal/app"
	"movie-app/metadata/internal/v1/controller/metadata"
	"movie-app/metadata/internal/v1/handler/http"
	"movie-app/metadata/internal/v1/repository"
)

func initApp(appType string) {
	repo := repository.New()
	ctrl := metadata.New(repo)

	if appType == "gin" {
		ginHandler := http.NewGinHandler(ctrl)
		err := app.NewGinApp(ginHandler, "3012").Run()
		if err != nil {
			log.Fatalln("failed to start gin app: ", err.Error())
		}
	}
}

func main() {
	log.Println("Starting movie metadata service")
	initApp("gin")
}
