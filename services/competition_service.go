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
