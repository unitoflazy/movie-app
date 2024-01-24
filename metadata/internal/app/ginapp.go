package app

import (
	"github.com/gin-gonic/gin"
	"movie-app/metadata/internal/v1/handler/http"
)

type GinApp struct {
	handler *http.GinHandler
	port    string
}

func NewGinApp(handler *http.GinHandler, port string) *GinApp {
	return &GinApp{
		handler: handler,
		port:    port,
	}
}

func (a *GinApp) Run() error {
	r := gin.Default()

	v1Router := r.Group("api/v1/metadata")
	v1Router.GET("", a.handler.GetMetadata)

	err := r.Run(":" + a.port)
	return err
}
