package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"movie-app/metadata/internal/v1/controller/metadata"
	"movie-app/metadata/internal/v1/repository"
	"net/http"
)

type GinHandler struct {
	ctrl *metadata.Controller
}

func NewGinHandler(ctrl *metadata.Controller) *GinHandler {
	return &GinHandler{ctrl}
}

func (h *GinHandler) GetMetadata(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	m, err := h.ctrl.Get(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		log.Printf("handler get error: %v\n", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, m)
}
