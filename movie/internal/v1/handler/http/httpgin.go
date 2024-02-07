package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"movie-app/movie/internal/v1/controller"
	"movie-app/movie/internal/v1/gateway"
	"net/http"
)

type GinHandler struct {
	ctrl *controller.Controller
}

func NewGinHandler(ctrl *controller.Controller) *GinHandler {
	return &GinHandler{ctrl: ctrl}
}

func (h *GinHandler) GetMovieDetails(ctx *gin.Context) {
	id := ctx.Param("id")
	details, err := h.ctrl.GetMovieDetails(ctx, id)
	if errors.Is(err, gateway.ErrNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": details,
	})
}
