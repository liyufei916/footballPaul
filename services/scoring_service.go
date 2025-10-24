package services

import (
	"github.com/liyufei916/footballPaul/database"
	"github.com/liyufei916/footballPaul/models"
	"github.com/liyufei916/footballPaul/utils"
	"gorm.io/gorm"
)

type ScoringService struct{}

func NewScoringService() *ScoringService {
	return &ScoringService{}
}

func (s *ScoringService) ScoreMatchPredictions(tx *gorm.DB, matchID uint, actualHome, actualAway int) (int, error) {
	var predictions []models.Prediction
	if err := tx.Where("match_id = ? AND is_scored = ?", matchID, false).Find(&predictions).Error; err != nil {
		return 0, err
	}

	scoredCount := 0
	for _, prediction := range predictions {
		points := utils.CalculatePoints(
			prediction.PredictedHomeScore,
			prediction.PredictedAwayScore,
			actualHome,
			actualAway,
		)

		if err := tx.Model(&prediction).Updates(map[string]interface{}{
			"points_earned": points,
			"is_scored":     true,
		}).Error; err != nil {
			return scoredCount, err
		}

		if err := tx.Model(&models.User{}).Where("id = ?", prediction.UserID).
			UpdateColumn("total_points", gorm.Expr("total_points + ?", points)).Error; err != nil {
			return scoredCount, err
		}

		scoredCount++
	}

	return scoredCount, nil
}

func (s *ScoringService) CalculatePredictionAccuracy(userID uint) (float64, error) {
	var totalPredictions int64
	var correctPredictions int64

	if err := database.DB.Model(&models.Prediction{}).
		Where("user_id = ? AND is_scored = ?", userID, true).
		Count(&totalPredictions).Error; err != nil {
		return 0, err
	}

	if totalPredictions == 0 {
		return 0, nil
	}

	if err := database.DB.Model(&models.Prediction{}).
		Where("user_id = ? AND is_scored = ? AND points_earned > ?", userID, true, 0).
		Count(&correctPredictions).Error; err != nil {
		return 0, err
	}

	return float64(correctPredictions) / float64(totalPredictions) * 100, nil
}
