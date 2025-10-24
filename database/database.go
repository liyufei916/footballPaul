package database

import (
	"fmt"
	"log"

	"github.com/liyufei916/footballPaul/config"
	"github.com/liyufei916/footballPaul/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase(cfg *config.Config) error {
	var err error
	dsn := cfg.Database.GetDSN()

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")

	if err := AutoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed")

	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Match{},
		&models.Prediction{},
		&models.ScoringRule{},
	)
}

func SeedDefaultScoringRules() error {
	var count int64
	if err := DB.Model(&models.ScoringRule{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("Scoring rules already exist, skipping seed")
		return nil
	}

	rules := []models.ScoringRule{
		{
			RuleName:    "完全正确",
			RuleType:    models.RuleTypeExactScore,
			Points:      10,
			Description: "比分完全一致",
		},
		{
			RuleName:    "猜中胜负+净胜球",
			RuleType:    models.RuleTypeGoalDifference,
			Points:      7,
			Description: "结果和净胜球都正确",
		},
		{
			RuleName:    "猜中胜负",
			RuleType:    models.RuleTypeCorrectWinner,
			Points:      5,
			Description: "只猜中胜/平/负",
		},
		{
			RuleName:    "猜中一方得分",
			RuleType:    models.RuleTypeOneScoreCorrect,
			Points:      3,
			Description: "猜中任一队伍得分",
		},
	}

	result := DB.Create(&rules)
	if result.Error != nil {
		return result.Error
	}

	log.Printf("Seeded %d scoring rules", len(rules))
	return nil
}
