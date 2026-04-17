package services

import (
	"errors"
	"time"

	"github.com/liyufei916/footballPaul/database"
	"github.com/liyufei916/footballPaul/models"
	"gorm.io/gorm"
)

type PredictionService struct {
	matchService *MatchService
}

func NewPredictionService() *PredictionService {
	return &PredictionService{
		matchService: NewMatchService(),
	}
}

func (s *PredictionService) CreatePrediction(userID, matchID uint, predictedHome, predictedAway int) (*models.Prediction, error) {
	deadlinePassed, err := s.matchService.IsDeadlinePassed(matchID)
	if err != nil {
		return nil, err
	}
	if deadlinePassed {
		return nil, errors.New("prediction deadline has passed")
	}

	var existingPrediction models.Prediction
	result := database.DB.Where("user_id = ? AND match_id = ?", userID, matchID).First(&existingPrediction)
	if result.Error == nil {
		return nil, errors.New("prediction already exists for this match")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	prediction := &models.Prediction{
		UserID:             userID,
		MatchID:            matchID,
		PredictedHomeScore: predictedHome,
		PredictedAwayScore: predictedAway,
		PredictedAt:        time.Now(),
		IsScored:           false,
	}

	if err := database.DB.Create(prediction).Error; err != nil {
		return nil, err
	}

	return prediction, nil
}

func (s *PredictionService) UpdatePrediction(predictionID, userID uint, predictedHome, predictedAway int) (*models.Prediction, error) {
	var prediction models.Prediction
	result := database.DB.Where("id = ? AND user_id = ?", predictionID, userID).First(&prediction)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("prediction not found")
		}
		return nil, result.Error
	}

	if prediction.IsScored {
		return nil, errors.New("cannot update scored prediction")
	}

	deadlinePassed, err := s.matchService.IsDeadlinePassed(prediction.MatchID)
	if err != nil {
		return nil, err
	}
	if deadlinePassed {
		return nil, errors.New("prediction deadline has passed")
	}

	prediction.PredictedHomeScore = predictedHome
	prediction.PredictedAwayScore = predictedAway

	if err := database.DB.Save(&prediction).Error; err != nil {
		return nil, err
	}

	return &prediction, nil
}

func (s *PredictionService) GetUserPredictions(userID uint, competitionID uint, includeMatch bool) ([]models.Prediction, error) {
	var predictions []models.Prediction
	query := database.DB.Where("user_id = ?", userID)

	if competitionID > 0 {
		query = query.Joins("JOIN matches ON predictions.match_id = matches.id").
			Where("matches.competition_id = ?", competitionID)
	}

	query = query.Order("predictions.created_at DESC")

	if includeMatch {
		query = query.Preload("Match").Preload("Match.Competition")
	}

	result := query.Find(&predictions)
	if result.Error != nil {
		return nil, result.Error
	}

	return predictions, nil
}

func (s *PredictionService) GetMatchPredictions(matchID uint, includeUser bool) ([]models.Prediction, error) {
	var predictions []models.Prediction
	query := database.DB.Where("match_id = ?", matchID)

	if includeUser {
		query = query.Preload("User")
	}

	result := query.Find(&predictions)
	if result.Error != nil {
		return nil, result.Error
	}

	return predictions, nil
}

// GetMatchPredictionsWithUsers 查询比赛的所有预测，并加载用户和比赛信息
func (s *PredictionService) GetMatchPredictionsWithUsers(matchID uint) ([]models.Prediction, error) {
	var predictions []models.Prediction
	result := database.DB.Where("match_id = ?", matchID).
		Preload("User").
		Preload("Match").
		Preload("Match.Competition").
		Find(&predictions)
	if result.Error != nil {
		return nil, result.Error
	}
	return predictions, nil
}

func (s *PredictionService) GetPredictionByID(id uint) (*models.Prediction, error) {
	var prediction models.Prediction
	result := database.DB.Preload("User").Preload("Match").Preload("Match.Competition").First(&prediction, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("prediction not found")
		}
		return nil, result.Error
	}
	return &prediction, nil
}
