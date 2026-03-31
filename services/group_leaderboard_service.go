package services

import (
	"github.com/liyufei916/footballPaul/database"
	"github.com/liyufei916/footballPaul/models"
	"gorm.io/gorm"
)

type GroupLeaderboardEntry struct {
	Rank             int   `json:"rank"`
	UserID           uint  `json:"user_id"`
	Username         string `json:"username"`
	TotalPoints      int   `json:"total_points"`
	PredictionsCount int64 `json:"predictions_count"`
	ExactScores      int   `json:"exact_scores"`
	CorrectWinners   int   `json:"correct_winners"`
}

type GroupLeaderboardService struct{}

func NewGroupLeaderboardService() *GroupLeaderboardService {
	return &GroupLeaderboardService{}
}

// GetGroupLeaderboard returns the leaderboard for a specific competition within a group
func (s *GroupLeaderboardService) GetGroupLeaderboard(groupID, competitionID uint, limit int) ([]GroupLeaderboardEntry, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	// Verify the competition is tracked by this group
	var gcCount int64
	database.DB.Model(&models.GroupCompetition{}).
		Where("group_id = ? AND competition_id = ?", groupID, competitionID).
		Count(&gcCount)
	if gcCount == 0 {
		return nil, ErrCompetitionNotFound
	}

	var results []struct {
		UserID           uint   `gorm:"column:user_id"`
		Username         string `gorm:"column:username"`
		TotalPoints      int    `gorm:"column:total_points"`
		PredictionsCount int64  `gorm:"column:predictions_count"`
		ExactScores      int    `gorm:"column:exact_scores"`
		CorrectWinners   int    `gorm:"column:correct_winners"`
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Table("users u").
			Select(`
				u.id as user_id,
				u.username as username,
				COALESCE(SUM(p.points_earned), 0) as total_points,
				COUNT(p.id) as predictions_count,
				COUNT(CASE WHEN p.points_earned = 10 THEN 1 END) as exact_scores,
				COUNT(CASE WHEN p.points_earned IN (5, 7) THEN 1 END) as correct_winners
			`).
			Joins("INNER JOIN group_members gm ON gm.user_id = u.id AND gm.group_id = ?", groupID).
			Joins("INNER JOIN predictions p ON p.user_id = u.id").
			Joins("INNER JOIN matches m ON m.id = p.match_id AND m.competition_id = ?", competitionID).
			Where("p.is_scored = ?", true).
			Group("u.id").
			Order("total_points DESC, exact_scores DESC").
			Limit(limit).
			Scan(&results).Error
	})

	if err != nil {
		return nil, err
	}

	entries := make([]GroupLeaderboardEntry, 0, len(results))
	for i, r := range results {
		entries = append(entries, GroupLeaderboardEntry{
			Rank:             i + 1,
			UserID:           r.UserID,
			Username:         r.Username,
			TotalPoints:      r.TotalPoints,
			PredictionsCount: r.PredictionsCount,
			ExactScores:      r.ExactScores,
			CorrectWinners:   r.CorrectWinners,
		})
	}

	return entries, nil
}

// GetUserRankInGroup returns a specific user's rank in a group's competition leaderboard
func (s *GroupLeaderboardService) GetUserRankInGroup(groupID, competitionID, userID uint) (*GroupLeaderboardEntry, error) {
	entries, err := s.GetGroupLeaderboard(groupID, competitionID, 100)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if e.UserID == userID {
			return &e, nil
		}
	}

	var username string
	database.DB.Model(&models.User{}).Where("id = ?", userID).Pluck("username", &username)

	return &GroupLeaderboardEntry{
		Rank:             0,
		UserID:           userID,
		Username:         username,
		TotalPoints:      0,
		PredictionsCount: 0,
		ExactScores:      0,
		CorrectWinners:   0,
	}, nil
}
