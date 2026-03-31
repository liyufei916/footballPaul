package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePoints_ExactScore(t *testing.T) {
	// Exact match: 10 points
	assert.Equal(t, PointsExactScore, CalculatePoints(2, 1, 2, 1))
	assert.Equal(t, PointsExactScore, CalculatePoints(0, 0, 0, 0))
	assert.Equal(t, PointsExactScore, CalculatePoints(1, 1, 1, 1))
}

func TestCalculatePoints_GoalDifference(t *testing.T) {
	// Same result + same goal diff: 7 points
	assert.Equal(t, PointsGoalDifference, CalculatePoints(2, 1, 3, 2)) // home win, diff=1
	assert.Equal(t, PointsGoalDifference, CalculatePoints(1, 2, 0, 1)) // away win, diff=1
	assert.Equal(t, PointsGoalDifference, CalculatePoints(0, 0, 3, 3)) // draw, diff=0
}

func TestCalculatePoints_CorrectWinner(t *testing.T) {
	// Same result + different goal diff: 5 points
	assert.Equal(t, PointsCorrectWinner, CalculatePoints(2, 1, 4, 1)) // home win, diff=1 vs 3
	assert.Equal(t, PointsCorrectWinner, CalculatePoints(1, 2, 0, 5)) // away win, diff=1 vs 5
	assert.Equal(t, PointsCorrectWinner, CalculatePoints(2, 1, 3, 1)) // home win, diff=1 vs 2
}

func TestCalculatePoints_Incorrect(t *testing.T) {
	// Wrong result + no score matches: 0 points
	assert.Equal(t, PointsIncorrect, CalculatePoints(2, 1, 0, 2)) // diff: home 2 vs 0, away 1 vs 2
}

func TestGetMatchResult(t *testing.T) {
	assert.Equal(t, HomeWin, getMatchResult(3, 1))
	assert.Equal(t, AwayWin, getMatchResult(0, 2))
	assert.Equal(t, Draw, getMatchResult(1, 1))
	assert.Equal(t, Draw, getMatchResult(0, 0))
}

func TestAbs(t *testing.T) {
	assert.Equal(t, 5, abs(5))
	assert.Equal(t, 5, abs(-5))
	assert.Equal(t, 0, abs(0))
}

func TestPointsConstants(t *testing.T) {
	assert.Equal(t, 10, PointsExactScore)
	assert.Equal(t, 7, PointsGoalDifference)
	assert.Equal(t, 5, PointsCorrectWinner)
	assert.Equal(t, 3, PointsOneScoreCorrect)
	assert.Equal(t, 0, PointsIncorrect)
}
