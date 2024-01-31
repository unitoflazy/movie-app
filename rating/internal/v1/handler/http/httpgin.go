package http

import (
	"github.com/gin-gonic/gin"
	"movie-app/rating/internal/v1/controller/rating"
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
	id := model.RecordID(ctx.Request.FormValue("id"))
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	ratingType := model.RecordType(ctx.Request.FormValue("type"))
	if ratingType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "type is required"})
		return
	}

	aggregatedRating, err := h.ctrl.GetAggregatedRating(ctx, ratingType, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": aggregatedRating})
}

func (h *GinHandler) PutRating(ctx *gin.Context) {
	var req model.RateMovieRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ctrl.PutRating(ctx, req.RecordType, req.RecordID, &model.Rating{
		Value:  req.Value,
		UserID: req.UserID,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
