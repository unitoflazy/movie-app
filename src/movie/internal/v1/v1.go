package v1

import (
	"github.com/gin-gonic/gin"
	"movie-app/movie/internal/v1/controller"
	metahttp "movie-app/movie/internal/v1/gateway/metadata/http"
	ratinghttp "movie-app/movie/internal/v1/gateway/rating/http"
	"movie-app/movie/internal/v1/handler/http"
	"movie-app/pkg/discovery"
)

type GinApp struct {
	port     string
	registry *discovery.Registry
}

func NewGinApp(port string, registry *discovery.Registry) *GinApp {
	return &GinApp{
		port:     port,
		registry: registry,
	}
}

func (g *GinApp) Run() error {
	metadataGW := metahttp.New(g.registry)
	ratingGW := ratinghttp.New(g.registry)
	ctrl := controller.New(ratingGW, metadataGW)
	handler := http.NewGinHandler(ctrl)

	r := gin.Default()
	v1Router := r.Group("api/v1/api")
	v1Router.GET("/:id", handler.GetMovieDetails)

	return r.Run(":" + g.port)
}
