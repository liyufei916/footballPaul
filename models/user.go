package models

import (
	"time"
)

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Username    string    `gorm:"unique;not null" json:"username"`
	Email       string    `gorm:"unique;not null" json:"email"`
	Password    string    `gorm:"not null" json:"-"`
	TotalPoints int       `gorm:"default:0" json:"total_points"`
	IsAdmin     bool      `gorm:"default:false" json:"is_admin"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Predictions []Prediction `gorm:"foreignKey:UserID" json:"predictions,omitempty"`
}

type UserResponse struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	TotalPoints int       `json:"total_points"`
	IsAdmin     bool      `json:"is_admin"`
	CreatedAt   time.Time `json:"created_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		TotalPoints: u.TotalPoints,
		IsAdmin:     u.IsAdmin,
		CreatedAt:   u.CreatedAt,
	}
}
