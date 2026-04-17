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
		*req.PredictedHomeScore,
		*req.PredictedAwayScore,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 重新查询以获取完整关联数据
	prediction, _ = h.predictionService.GetPredictionByID(prediction.ID)

	c.JSON(http.StatusCreated, gin.H{
		"success":    true,
		"prediction": prediction,
		"message":    "预测提交成功",
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
		*req.PredictedHomeScore,
		*req.PredictedAwayScore,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证请求中的 match_id 与预测记录一致（防止误传数据）
	if req.MatchID > 0 && prediction.MatchID != req.MatchID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "match_id mismatch"})
		return
	}

	// 重新查询以获取完整关联数据（User、Match、Competition）
	prediction, _ = h.predictionService.GetPredictionByID(prediction.ID)

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

	competitionIDStr := c.Query("competition_id")
	competitionID, _ := strconv.ParseUint(competitionIDStr, 10, 32)

	predictions, err := h.predictionService.GetUserPredictions(userID.(uint), uint(competitionID), true)
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
	idStr := c.Param("id")
	matchID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}

	// 需要同时加载 User 和 Match+Competition 信息
	predictions, err := h.predictionService.GetMatchPredictionsWithUsers(uint(matchID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"predictions": predictions,
		"count":       len(predictions),
	})
}
