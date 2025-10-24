package services

import (
	"errors"
	"time"

	"github.com/liyufei916/footballPaul/database"
	"github.com/liyufei916/footballPaul/models"
	"gorm.io/gorm"
)

type MatchService struct{}

func NewMatchService() *MatchService {
	return &MatchService{}
}

func (s *MatchService) CreateMatch(homeTeam, awayTeam string, matchDate, deadline time.Time) (*models.Match, error) {
	match := &models.Match{
		HomeTeam:  homeTeam,
		AwayTeam:  awayTeam,
		MatchDate: matchDate,
		Deadline:  deadline,
		Status:    models.MatchStatusPending,
	}

	result := database.DB.Create(match)
	if result.Error != nil {
		return nil, result.Error
	}

	return match, nil
}

func (s *MatchService) GetMatchByID(id uint) (*models.Match, error) {
	var match models.Match
	result := database.DB.First(&match, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("match not found")
		}
		return nil, result.Error
	}
	return &match, nil
}

func (s *MatchService) GetMatches(status models.MatchStatus, limit int) ([]models.Match, error) {
	var matches []models.Match
	query := database.DB

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Order("match_date ASC").Find(&matches)
	if result.Error != nil {
		return nil, result.Error
	}

	return matches, nil
}

func (s *MatchService) UpdateMatchResult(matchID uint, homeScore, awayScore int) error {
	match, err := s.GetMatchByID(matchID)
	if err != nil {
		return err
	}

	if match.Status == models.MatchStatusFinished {
		return errors.New("match already finished")
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updates := map[string]interface{}{
		"home_score": homeScore,
		"away_score": awayScore,
		"status":     models.MatchStatusFinished,
	}

	if err := tx.Model(&models.Match{}).Where("id = ?", matchID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	scoringService := NewScoringService()
	scoredCount, err := scoringService.ScoreMatchPredictions(tx, matchID, homeScore, awayScore)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *MatchService) UpdateMatchStatus(matchID uint, status models.MatchStatus) error {
	result := database.DB.Model(&models.Match{}).Where("id = ?", matchID).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("match not found")
	}
	return nil
}

func (s *MatchService) IsDeadlinePassed(matchID uint) (bool, error) {
	match, err := s.GetMatchByID(matchID)
	if err != nil {
		return false, err
	}
	return time.Now().After(match.Deadline), nil
}
