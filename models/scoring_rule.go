package models

import (
	"time"
)

type RuleType string

const (
	RuleTypeExactScore      RuleType = "exact_score"
	RuleTypeCorrectWinner   RuleType = "correct_winner"
	RuleTypeGoalDifference  RuleType = "goal_difference"
	RuleTypeOneScoreCorrect RuleType = "one_score_correct"
)

type ScoringRule struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RuleName    string    `gorm:"not null" json:"rule_name"`
	RuleType    RuleType  `gorm:"not null" json:"rule_type"`
	Points      int       `gorm:"not null" json:"points"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
