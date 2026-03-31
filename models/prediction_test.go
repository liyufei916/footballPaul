package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPredictionRequest_Fields(t *testing.T) {
	req := PredictionRequest{
		MatchID:            10,
		PredictedHomeScore: 2,
		PredictedAwayScore: 1,
	}

	assert.Equal(t, uint(10), req.MatchID)
	assert.Equal(t, 2, req.PredictedHomeScore)
	assert.Equal(t, 1, req.PredictedAwayScore)
}

func TestPredictionResponse_Fields(t *testing.T) {
	homeScore := 3
	awayScore := 1
	resp := PredictionResponse{
		ID:                 5,
		UserID:             1,
		MatchID:            10,
		PredictedHomeScore: 2,
		PredictedAwayScore: 1,
		PointsEarned:       10,
		IsScored:           true,
		PredictedAt:        time.Now(),
		HomeTeam:           "Team A",
		AwayTeam:           "Team B",
		ActualHomeScore:    &homeScore,
		ActualAwayScore:    &awayScore,
	}

	assert.Equal(t, uint(5), resp.ID)
	assert.Equal(t, 2, resp.PredictedHomeScore)
	assert.True(t, resp.IsScored)
	assert.Equal(t, "Team A", resp.HomeTeam)
	assert.NotNil(t, resp.ActualHomeScore)
	assert.Equal(t, 3, *resp.ActualHomeScore)
}

func TestPrediction_DefaultValues(t *testing.T) {
	p := Prediction{}
	assert.Equal(t, 0, p.PointsEarned)
	assert.False(t, p.IsScored)
}
