package v1

import (
	"github.com/gin-gonic/gin"
	"movie-app/rating/internal/v1/controller/rating"
	"movie-app/rating/internal/v1/handler/http"
	"movie-app/rating/internal/v1/repository/memory"
)

type GinApp struct {
	port string
}

func NewGinApp(port string) *GinApp {
	return &GinApp{
		port: port,
	}
}

func (g *GinApp) Run() error {
	repo := memory.New()
	ctrl := rating.New(repo)
	handler := http.NewGinHandler(ctrl)

	r := gin.Default()
	v1Router := r.Group("api/v1/rating")

	v1Router.GET("/:id", handler.GetAggregatedRating)
	v1Router.PUT("/:id", handler.PutRating)

	return r.Run(":" + g.port)
}
