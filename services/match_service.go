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

func (s *MatchService) CreateMatch(competitionID uint, homeTeam, awayTeam string, matchDate, deadline time.Time) (*models.Match, error) {
	match := &models.Match{
		CompetitionID: competitionID,
		HomeTeam:      homeTeam,
		AwayTeam:      awayTeam,
		MatchDate:     matchDate,
		Deadline:      deadline,
		Status:        models.MatchStatusPending,
	}

	result := database.DB.Create(match)
	if result.Error != nil {
		return nil, result.Error
	}

	return match, nil
}

func (s *MatchService) GetMatchByID(id uint) (*models.Match, error) {
	var match models.Match
	result := database.DB.Preload("Competition").First(&match, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("match not found")
		}
		return nil, result.Error
	}
	// 动态计算状态
	match.Status = match.EffectiveStatus()
	return &match, nil
}

func (s *MatchService) GetMatches(status models.MatchStatus, competitionID uint, limit int) ([]models.Match, error) {
	var matches []models.Match
	query := database.DB.Preload("Competition")

	// 根据动态状态过滤：
	// - pending: match_date > now AND 未录比分
	// - ongoing: match_date <= now AND 未录比分
	// - finished: 已录比分（home_score IS NOT NULL）
	now := time.Now()
	switch status {
	case models.MatchStatusPending:
		query = query.Where("match_date > ? AND home_score IS NULL", now)
	case models.MatchStatusOngoing:
		query = query.Where("match_date <= ? AND home_score IS NULL", now)
	case models.MatchStatusFinished:
		query = query.Where("home_score IS NOT NULL")
	}

	if competitionID > 0 {
		query = query.Where("competition_id = ?", competitionID)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Order("match_date ASC").Find(&matches)
	if result.Error != nil {
		return nil, result.Error
	}

	// 动态计算每个比赛的状态
	for i := range matches {
		matches[i].Status = matches[i].EffectiveStatus()
	}

	return matches, nil
}

func (s *MatchService) UpdateMatchResult(matchID uint, homeScore, awayScore int) error {
	match, err := s.GetMatchByID(matchID)
	if err != nil {
		return err
	}

	// 必须比赛已开始才能录入结果（不允许赛前录入假结果）
	if !match.HasStarted() {
		return errors.New("比赛尚未开始，无法录入结果")
	}

	// 已结束的比赛不允许重复录入
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
	_, err = scoringService.ScoreMatchPredictions(tx, matchID, homeScore, awayScore)
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
