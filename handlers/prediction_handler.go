package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liyufei916/footballPaul/models"
	"github.com/liyufei916/footballPaul/services"
)

type PredictionHandler struct {
	predictionService *services.PredictionService
}

func NewPredictionHandler() *PredictionHandler {
	return &PredictionHandler{
		predictionService: services.NewPredictionService(),
	}
}

func (h *PredictionHandler) CreatePrediction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.PredictionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prediction, err := h.predictionService.CreatePrediction(
		userID.(uint),
		req.MatchID,
		req.PredictedHomeScore,
		req.PredictedAwayScore,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":       true,
		"prediction_id": prediction.ID,
		"message":       "预测提交成功",
	})
}

func (h *PredictionHandler) UpdatePrediction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid prediction id"})
		return
	}

	var req models.PredictionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prediction, err := h.predictionService.UpdatePrediction(
		uint(id),
		userID.(uint),
		req.PredictedHomeScore,
		req.PredictedAwayScore,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"prediction": prediction,
		"message":    "预测更新成功",
	})
}

func (h *PredictionHandler) GetUserPredictions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	predictions, err := h.predictionService.GetUserPredictions(userID.(uint), true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"predictions": predictions,
		"count":       len(predictions),
	})
}

func (h *PredictionHandler) GetMatchPredictions(c *gin.Context) {
	idStr := c.Param("matchId")
	matchID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}

	predictions, err := h.predictionService.GetMatchPredictions(uint(matchID), false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"predictions": predictions,
		"count":       len(predictions),
	})
}
