package models

import (
	"time"

	"gorm.io/gorm"
)

type MatchStatus string

const (
	MatchStatusPending  MatchStatus = "pending"
	MatchStatusOngoing MatchStatus = "ongoing"
	MatchStatusFinished MatchStatus = "finished"
)

type Match struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	CompetitionID uint        `gorm:"index" json:"competition_id"`
	HomeTeam      string      `gorm:"not null" json:"home_team"`
	AwayTeam      string      `gorm:"not null" json:"away_team"`
	MatchDate     time.Time   `gorm:"not null" json:"match_date"`
	HomeScore     *int        `json:"home_score"`
	AwayScore     *int        `json:"away_score"`
	Status        MatchStatus `gorm:"default:'pending'" json:"status"`
	Deadline      time.Time   `gorm:"not null" json:"deadline"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`

	Competition  Competition  `gorm:"foreignKey:CompetitionID" json:"competition,omitempty"`
	Predictions []Prediction `gorm:"foreignKey:MatchID" json:"predictions,omitempty"`
}

type MatchResult struct {
	HomeScore int `json:"home_score" binding:"required,min=0"`
	AwayScore int `json:"away_score" binding:"required,min=0"`
}

// EffectiveStatus 根据当前时间和比分录入情况，动态计算比赛状态
// - 已录入比分 -> finished
// - 比赛时间已到但未录比分 -> ongoing
// - 比赛时间未到 -> pending
func (m *Match) EffectiveStatus() MatchStatus {
	if m.HomeScore != nil && m.AwayScore != nil {
		return MatchStatusFinished
	}
	if time.Now().After(m.MatchDate) {
		return MatchStatusOngoing
	}
	return MatchStatusPending
}

// HasStarted 检查比赛是否已到开始时间
func (m *Match) HasStarted() bool {
	return time.Now().After(m.MatchDate)
}

// AfterCreate 创建比赛后自动更新赛事的 match_count
func (m *Match) AfterCreate(tx *gorm.DB) error {
	return tx.Model(&Competition{}).Where("id = ?", m.CompetitionID).
		UpdateColumn("match_count", gorm.Expr("match_count + 1")).Error
}
