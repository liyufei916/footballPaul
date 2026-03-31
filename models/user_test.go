package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserToResponse(t *testing.T) {
	now := time.Now()
	user := User{
		ID:          42,
		Username:    "testuser",
		Email:       "test@example.com",
		Password:    "secret-should-not-leak",
		TotalPoints: 150,
		IsAdmin:     true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	resp := user.ToResponse()

	assert.Equal(t, uint(42), resp.ID)
	assert.Equal(t, "testuser", resp.Username)
	assert.Equal(t, "test@example.com", resp.Email)
	assert.Equal(t, 150, resp.TotalPoints)
	assert.True(t, resp.IsAdmin)
	assert.Equal(t, now, resp.CreatedAt)
}

func TestUserToResponse_NonAdmin(t *testing.T) {
	user := User{IsAdmin: false}
	resp := user.ToResponse()
	assert.False(t, resp.IsAdmin)
}

func TestUserToResponse_Zeros(t *testing.T) {
	user := User{}
	resp := user.ToResponse()
	assert.Equal(t, uint(0), resp.ID)
	assert.Equal(t, "", resp.Username)
	assert.Equal(t, 0, resp.TotalPoints)
}
