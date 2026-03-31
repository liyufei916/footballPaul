package models

import (
	"time"
)

// Group represents a prediction group
type Group struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:50;not null" json:"name"`
	InviteCode string    `gorm:"size:6;uniqueIndex;not null" json:"invite_code"`
	OwnerID    uint      `gorm:"not null;index" json:"owner_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Owner         User             `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Members       []GroupMember    `gorm:"foreignKey:GroupID" json:"members,omitempty"`
	Competitions []GroupCompetition `gorm:"foreignKey:GroupID" json:"competitions,omitempty"`
}

// GroupMember represents a user's membership in a group
type GroupMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GroupID   uint      `gorm:"not null;index;uniqueIndex:idx_group_user" json:"group_id"`
	UserID    uint      `gorm:"not null;index;uniqueIndex:idx_group_user" json:"user_id"`
	Role      string    `gorm:"type:text;default:'member'" json:"role"` // "admin" or "member"
	JoinedAt  time.Time `json:"joined_at"`

	Group Group `gorm:"foreignKey:GroupID" json:"-"`
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// GroupCompetition represents a competition tracked by a group
type GroupCompetition struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	GroupID       uint      `gorm:"not null;index;uniqueIndex:idx_group_competition" json:"group_id"`
	CompetitionID uint      `gorm:"not null;index;uniqueIndex:idx_group_competition" json:"competition_id"`
	CreatedAt     time.Time `json:"created_at"`

	Group       Group       `gorm:"foreignKey:GroupID" json:"-"`
	Competition Competition `gorm:"foreignKey:CompetitionID" json:"competition,omitempty"`
}

// GroupMemberResponse is the API response for a group member
type GroupMemberResponse struct {
	UserID   uint      `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

// GroupResponse is the API response for a group
type GroupResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	InviteCode    string    `json:"invite_code,omitempty"`
	OwnerID       uint      `json:"owner_id"`
	MemberCount   int       `json:"member_count,omitempty"`
	CompetitionCount int    `json:"competition_count,omitempty"`
	Role          string    `json:"role,omitempty"`
	JoinedAt      time.Time `json:"joined_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// ToResponse converts a Group to GroupResponse
func (g *Group) ToResponse() GroupResponse {
	return GroupResponse{
		ID:         g.ID,
		Name:       g.Name,
		InviteCode: g.InviteCode,
		OwnerID:    g.OwnerID,
		CreatedAt:  g.CreatedAt,
	}
}
