package models

import (
	"time"
)

type Prediction struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	UserID             uint      `gorm:"not null;index" json:"user_id"`
	MatchID            uint      `gorm:"not null;index" json:"match_id"`
	PredictedHomeScore int       `gorm:"not null" json:"predicted_home_score"`
	PredictedAwayScore int       `gorm:"not null" json:"predicted_away_score"`
	PointsEarned       int       `gorm:"default:0" json:"points_earned"`
	IsScored           bool      `gorm:"default:false" json:"is_scored"`
	PredictedAt        time.Time `json:"predicted_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Match Match `gorm:"foreignKey:MatchID" json:"match,omitempty"`
}

type PredictionRequest struct {
	MatchID            uint `json:"match_id" binding:"required"`
	PredictedHomeScore int  `json:"predicted_home_score" binding:"required,min=0"`
	PredictedAwayScore int  `json:"predicted_away_score" binding:"required,min=0"`
}

type PredictionResponse struct {
	ID                 uint      `json:"id"`
	UserID             uint      `json:"user_id"`
	MatchID            uint      `json:"match_id"`
	PredictedHomeScore int       `json:"predicted_home_score"`
	PredictedAwayScore int       `json:"predicted_away_score"`
	PointsEarned       int       `json:"points_earned"`
	IsScored           bool      `json:"is_scored"`
	PredictedAt        time.Time `json:"predicted_at"`
	HomeTeam           string    `json:"home_team,omitempty"`
	AwayTeam           string    `json:"away_team,omitempty"`
	ActualHomeScore    *int      `json:"actual_home_score,omitempty"`
	ActualAwayScore    *int      `json:"actual_away_score,omitempty"`
}
