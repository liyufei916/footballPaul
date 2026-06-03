package services

import (
	"errors"

	"github.com/liyufei916/footballPaul/database"
	"github.com/liyufei916/footballPaul/models"
	"gorm.io/gorm"
)

type CompetitionService struct{}

func NewCompetitionService() *CompetitionService {
	return &CompetitionService{}
}

func (s *CompetitionService) GetAllCompetitions() ([]models.Competition, error) {
	var competitions []models.Competition
	result := database.DB.Order("name ASC").Find(&competitions)
	if result.Error != nil {
		return nil, result.Error
	}
	return competitions, nil
}

func (s *CompetitionService) GetCompetitionByID(id uint) (*models.Competition, error) {
	var competition models.Competition
	result := database.DB.First(&competition, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("competition not found")
		}
		return nil, result.Error
	}
	return &competition, nil
}

func (s *CompetitionService) CreateCompetition(name, code, logo string) (*models.Competition, error) {
	competition := &models.Competition{
		Name: name,
		Code: code,
		Logo: logo,
	}
	result := database.DB.Create(competition)
	if result.Error != nil {
		return nil, result.Error
	}
	return competition, nil
}

func (s *CompetitionService) DeleteCompetition(competitionID uint) error {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Verify competition exists
	var competition models.Competition
	if err := tx.First(&competition, competitionID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("赛事不存在")
		}
		return err
	}

	// Get all match IDs of this competition
	var matchIDs []uint
	if err := tx.Model(&models.Match{}).Where("competition_id = ?", competitionID).Pluck("id", &matchIDs).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete predictions for all matches of this competition
	if len(matchIDs) > 0 {
		if err := tx.Where("match_id IN ?", matchIDs).Delete(&models.Prediction{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Delete all matches of this competition
	if err := tx.Where("competition_id = ?", competitionID).Delete(&models.Match{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete GroupCompetition entries for this competition
	if err := tx.Where("competition_id = ?", competitionID).Delete(&models.GroupCompetition{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the competition
	if err := tx.Delete(&models.Competition{}, competitionID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
