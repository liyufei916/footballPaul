package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liyufei916/footballPaul/models"
	"github.com/liyufei916/footballPaul/services"
)

type MatchHandler struct {
	matchService *services.MatchService
}

func NewMatchHandler() *MatchHandler {
	return &MatchHandler{
		matchService: services.NewMatchService(),
	}
}

type CreateMatchRequest struct {
	CompetitionID uint      `json:"competition_id" binding:"required"`
	HomeTeam     string    `json:"home_team" binding:"required"`
	AwayTeam     string    `json:"away_team" binding:"required"`
	MatchDate    time.Time `json:"match_date" binding:"required"`
	Deadline     time.Time `json:"deadline" binding:"required"`
}

func (h *MatchHandler) CreateMatch(c *gin.Context) {
	var req CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	match, err := h.matchService.CreateMatch(req.CompetitionID, req.HomeTeam, req.AwayTeam, req.MatchDate, req.Deadline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"match":   match,
	})
}

func (h *MatchHandler) GetMatch(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}

	match, err := h.matchService.GetMatchByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, match)
}

func (h *MatchHandler) GetMatches(c *gin.Context) {
	status := models.MatchStatus(c.Query("status"))
	competitionIDStr := c.Query("competition_id")
	limitStr := c.DefaultQuery("limit", "10")

	competitionID, _ := strconv.ParseUint(competitionIDStr, 10, 32)

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = 10
	}

	matches, err := h.matchService.GetMatches(status, uint(competitionID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"matches": matches,
		"count":   len(matches),
	})
}

func (h *MatchHandler) UpdateMatchResult(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}

	var req models.MatchResult
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.matchService.UpdateMatchResult(uint(id), req.HomeScore, req.AwayScore); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "比分已录入，评分完成",
	})
}
