package services

import (
	"github.com/liyufei916/footballPaul/database"
	"github.com/liyufei916/footballPaul/models"
)

type LeaderboardEntry struct {
	Rank             int    `json:"rank"`
	UserID           uint   `json:"user_id"`
	Username         string `json:"username"`
	TotalPoints      int    `json:"total_points"`
	PredictionsCount int64  `json:"predictions_count"`
}

type LeaderboardService struct{}

func NewLeaderboardService() *LeaderboardService {
	return &LeaderboardService{}
}

func (s *LeaderboardService) GetLeaderboard(limit int) ([]LeaderboardEntry, error) {
	var users []models.User

	query := database.DB.Order("total_points DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	leaderboard := make([]LeaderboardEntry, 0, len(users))
	for i, user := range users {
		var predictionCount int64
		database.DB.Model(&models.Prediction{}).Where("user_id = ?", user.ID).Count(&predictionCount)

		leaderboard = append(leaderboard, LeaderboardEntry{
			Rank:             i + 1,
			UserID:           user.ID,
			Username:         user.Username,
			TotalPoints:      user.TotalPoints,
			PredictionsCount: predictionCount,
		})
	}

	return leaderboard, nil
}

func (s *LeaderboardService) GetUserRank(userID uint) (int, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return 0, err
	}

	var rank int64
	if err := database.DB.Model(&models.User{}).
		Where("total_points > ?", user.TotalPoints).
		Count(&rank).Error; err != nil {
		return 0, err
	}

	return int(rank) + 1, nil
}
