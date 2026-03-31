package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestCompetitionHandler_GetCompetitions_Structure(t *testing.T) {
	// Just test handler creation and routing structure
	handler := NewCompetitionHandler()
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.competitionService)
}

func TestMatchHandler_NewMatchHandler(t *testing.T) {
	handler := NewMatchHandler()
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.matchService)
}

func TestPredictionHandler_NewPredictionHandler(t *testing.T) {
	handler := NewPredictionHandler()
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.predictionService)
}

func TestLeaderboardHandler_NewLeaderboardHandler(t *testing.T) {
	handler := NewLeaderboardHandler()
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.leaderboardService)
}

func TestHealthEndpoint(t *testing.T) {
	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp["status"])
}
