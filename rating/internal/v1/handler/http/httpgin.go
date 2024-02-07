package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"movie-app/rating/internal/v1/controller/rating"
	"movie-app/rating/internal/v1/repository"
	"movie-app/rating/pkg/model"
	"net/http"
)

type GinHandler struct {
	ctrl *rating.Controller
}

func NewGinHandler(ctrl *rating.Controller) *GinHandler {
	return &GinHandler{ctrl: ctrl}
}

func (h *GinHandler) GetAggregatedRating(ctx *gin.Context) {
	id := model.RecordID(ctx.Param("id"))
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	ratingType := model.RecordType(ctx.Query("type"))
	if ratingType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "type is required"})
		return
	}

	aggregatedRating, err := h.ctrl.GetAggregatedRating(ctx, ratingType, id)
	if errors.Is(err, repository.ErrNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, aggregatedRating)
}

func (h *GinHandler) PutRating(ctx *gin.Context) {
	id := model.RecordID(ctx.Param("id"))
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	ratingType := model.RecordType(ctx.Query("type"))
	if ratingType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "type is required"})
		return
	}

	var req model.RateMovieRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ctrl.PutRating(ctx, ratingType, id, &model.Rating{
		Value:  req.Value,
		UserID: req.UserID,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, true)
}
