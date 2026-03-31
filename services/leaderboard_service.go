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

func (s *LeaderboardService) GetLeaderboard(competitionID uint, limit int) ([]LeaderboardEntry, error) {
	type userWithStats struct {
		UserID           uint  `gorm:"column:id"`
		Username         string `gorm:"column:username"`
		TotalPoints      int    `gorm:"column:total_points"`
		PredictionsCount  int64  `gorm:"column:predictions_count"`
	}

	var results []userWithStats

	if competitionID > 0 {
		// Competition-specific leaderboard
		subQuery := database.DB.Model(&models.Prediction{}).
			Select("predictions.user_id, COUNT(*) as predictions_count, COALESCE(SUM(predictions.points_earned), 0) as total_points").
			Joins("JOIN matches ON predictions.match_id = matches.id").
			Where("matches.competition_id = ? AND predictions.is_scored = ?", competitionID, true).
			Group("predictions.user_id")

		query := database.DB.Table("users").
			Select("users.id, users.username, COALESCE(sq.total_points, 0) as total_points, COALESCE(sq.predictions_count, 0) as predictions_count").
			Joins("LEFT JOIN (?) as sq ON users.id = sq.user_id", subQuery).
			Where("sq.predictions_count > 0").
			Order("sq.total_points DESC, users.id ASC")

		if limit > 0 {
			query = query.Limit(limit)
		}

		if err := query.Scan(&results).Error; err != nil {
			return nil, err
		}
	} else {
		// Global leaderboard (all predictions)
		subQuery := database.DB.Model(&models.Prediction{}).
			Select("user_id, COUNT(*) as predictions_count, COALESCE(SUM(points_earned), 0) as total_points").
			Where("is_scored = ?", true).
			Group("user_id")

		query := database.DB.Table("users").
			Select("users.id, users.username, COALESCE(sq.total_points, 0) as total_points, COALESCE(sq.predictions_count, 0) as predictions_count").
			Joins("LEFT JOIN (?) as sq ON users.id = sq.user_id", subQuery).
			Where("sq.predictions_count > 0").
			Order("sq.total_points DESC, users.id ASC")

		if limit > 0 {
			query = query.Limit(limit)
		}

		if err := query.Scan(&results).Error; err != nil {
			return nil, err
		}
	}

	leaderboard := make([]LeaderboardEntry, 0, len(results))
	for i, r := range results {
		leaderboard = append(leaderboard, LeaderboardEntry{
			Rank:             i + 1,
			UserID:           r.UserID,
			Username:         r.Username,
			TotalPoints:      r.TotalPoints,
			PredictionsCount: r.PredictionsCount,
		})
	}

	return leaderboard, nil
}

func (s *LeaderboardService) GetUserRank(userID uint, competitionID uint) (int, error) {
	type userWithStats struct {
		UserID      uint `gorm:"column:id"`
		TotalPoints int  `gorm:"column:total_points"`
	}

	var currentUser userWithStats

	if competitionID > 0 {
		// Competition-specific rank
		subQuery := database.DB.Model(&models.Prediction{}).
			Select("user_id, COALESCE(SUM(points_earned), 0) as total_points").
			Joins("JOIN matches ON predictions.match_id = matches.id").
			Where("matches.competition_id = ? AND predictions.is_scored = ?", competitionID, true).
			Group("predictions.user_id").
			Having("user_id = ?", userID)

		if err := database.DB.Table("users").
			Select("users.id, COALESCE(sq.total_points, 0) as total_points").
			Joins("LEFT JOIN (?) as sq ON users.id = sq.user_id", subQuery).
			Where("users.id = ?", userID).
			Scan(&currentUser).Error; err != nil {
			return 0, err
		}

		// Count users with higher points in this competition
		var rank int64
		subQueryHigher := database.DB.Model(&models.Prediction{}).
			Select("user_id, COALESCE(SUM(points_earned), 0) as total_points").
			Joins("JOIN matches ON predictions.match_id = matches.id").
			Where("matches.competition_id = ? AND predictions.is_scored = ?", competitionID, true).
			Group("predictions.user_id").
			Having("COALESCE(SUM(points_earned), 0) > ?", currentUser.TotalPoints)

		if err := database.DB.Model(&models.User{}).
			Joins("LEFT JOIN (?) as sq ON users.id = sq.user_id", subQueryHigher).
			Where("sq.user_id IS NOT NULL").
			Count(&rank).Error; err != nil {
			return 0, err
		}

		return int(rank) + 1, nil
	} else {
		// Global rank
		if err := database.DB.Table("users").
			Select("users.id, COALESCE(sq.total_points, 0) as total_points").
			Joins("LEFT JOIN (SELECT user_id, COALESCE(SUM(points_earned), 0) as total_points FROM predictions WHERE is_scored = true GROUP BY user_id) as sq ON users.id = sq.user_id").
			Where("users.id = ?", userID).
			Scan(&currentUser).Error; err != nil {
			return 0, err
		}

		var rank int64
		if err := database.DB.Model(&models.User{}).
			Joins("LEFT JOIN (SELECT user_id, COALESCE(SUM(points_earned), 0) as total_points FROM predictions WHERE is_scored = true GROUP BY user_id) as sq ON users.id = sq.user_id").
			Where("COALESCE(sq.total_points, 0) > ?", currentUser.TotalPoints).
			Count(&rank).Error; err != nil {
			return 0, err
		}

		return int(rank) + 1, nil
	}
}

func (s *LeaderboardService) GetUserPointsByCompetition(userID uint) (map[uint]int, error) {
	type competitionPoints struct {
		CompetitionID uint
		TotalPoints   int
	}

	var results []competitionPoints

	if err := database.DB.Model(&models.Prediction{}).
		Select("matches.competition_id, COALESCE(SUM(predictions.points_earned), 0) as total_points").
		Joins("JOIN matches ON predictions.match_id = matches.id").
		Where("predictions.user_id = ? AND predictions.is_scored = ?", userID, true).
		Group("matches.competition_id").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	pointsMap := make(map[uint]int)
	for _, r := range results {
		pointsMap[r.CompetitionID] = r.TotalPoints
	}

	return pointsMap, nil
}
