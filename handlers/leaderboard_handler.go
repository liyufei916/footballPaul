package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liyufei916/footballPaul/services"
)

type LeaderboardHandler struct {
	leaderboardService *services.LeaderboardService
}

func NewLeaderboardHandler() *LeaderboardHandler {
	return &LeaderboardHandler{
		leaderboardService: services.NewLeaderboardService(),
	}
}

func (h *LeaderboardHandler) GetLeaderboard(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	competitionIDStr := c.Query("competition_id")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = 50
	}

	competitionID, _ := strconv.ParseUint(competitionIDStr, 10, 32)

	leaderboard, err := h.leaderboardService.GetLeaderboard(uint(competitionID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rankings": leaderboard,
	})
}

func (h *LeaderboardHandler) GetUserRank(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	competitionIDStr := c.Query("competition_id")
	competitionID, _ := strconv.ParseUint(competitionIDStr, 10, 32)

	rank, err := h.leaderboardService.GetUserRank(userID.(uint), uint(competitionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rank": rank,
	})
}
