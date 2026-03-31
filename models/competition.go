package models

import (
	"time"
)

type Competition struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	Code      string    `gorm:"unique;not null" json:"code"`
	Logo      string    `json:"logo,omitempty"`
	MatchCount int      `gorm:"default:0" json:"match_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Matches []Match `gorm:"foreignKey:CompetitionID" json:"matches,omitempty"`
}
