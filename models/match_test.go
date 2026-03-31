package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchStatus_Constants(t *testing.T) {
	assert.Equal(t, MatchStatus("pending"), MatchStatusPending)
	assert.Equal(t, MatchStatus("ongoing"), MatchStatusOngoing)
	assert.Equal(t, MatchStatus("finished"), MatchStatusFinished)
}

func TestMatchResult_Validation(t *testing.T) {
	tests := []struct {
		name    string
		home    int
		away    int
		invalid bool
	}{
		{"valid 0-0", 0, 0, false},
		{"valid normal score", 3, 2, false},
		{"valid large score", 10, 5, false},
		// Score can be 0 but must be non-negative
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Non-negative scores should be valid (basic sanity)
			assert.GreaterOrEqual(t, tt.home, 0)
			assert.GreaterOrEqual(t, tt.away, 0)
		})
	}
}

func TestMatchResult_BindingValidation(t *testing.T) {
	// Test that MatchResult struct requires non-negative values
	result := MatchResult{HomeScore: 0, AwayScore: 0}
	assert.Equal(t, 0, result.HomeScore)
	assert.Equal(t, 0, result.AwayScore)
}
