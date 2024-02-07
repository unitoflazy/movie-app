package v1

import (
	"github.com/gin-gonic/gin"
	"movie-app/movie/internal/v1/controller"
	metahttp "movie-app/movie/internal/v1/gateway/metadata/http"
	ratinghttp "movie-app/movie/internal/v1/gateway/rating/http"
	"movie-app/movie/internal/v1/handler/http"
)

type GinApp struct {
	port        string
	metadataAdr string
	ratingAdr   string
}

func NewGinApp(port string, metadataAdr string, ratingAdr string) *GinApp {
	return &GinApp{
		port:        port,
		metadataAdr: metadataAdr,
		ratingAdr:   ratingAdr,
	}
}

func (g *GinApp) Run() error {
	metadataGW := metahttp.New(g.metadataAdr)
	ratingGW := ratinghttp.New(g.ratingAdr)
	ctrl := controller.New(ratingGW, metadataGW)
	handler := http.NewGinHandler(ctrl)

	r := gin.Default()
	v1Router := r.Group("api/v1/movie")
	v1Router.GET("/:id", handler.GetMovieDetails)

	return r.Run(":" + g.port)
}
