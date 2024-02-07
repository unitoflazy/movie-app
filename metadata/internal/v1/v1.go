package v1

import (
	"github.com/gin-gonic/gin"
	"movie-app/metadata/internal/v1/controller/metadata"
	"movie-app/metadata/internal/v1/handler/http"
	"movie-app/metadata/internal/v1/repository/memory"
)

type GinApp struct {
	port string
}

func NewGinApp(port string) *GinApp {
	return &GinApp{
		port: port,
	}
}

func (a *GinApp) Run() error {
	repo := memory.New()
	ctrl := metadata.New(repo)
	ginHandler := http.NewGinHandler(ctrl)

	r := gin.Default()

	v1Router := r.Group("api/v1/metadata")
	v1Router.GET("/:id", ginHandler.GetMetadata)

	return r.Run(":" + a.port)
}
