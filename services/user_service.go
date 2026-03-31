package services

import (
	"errors"

	"github.com/liyufei916/footballPaul/database"
	"github.com/liyufei916/footballPaul/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(username, email, password string) (*models.User, error) {
	// Check for existing user by email (case-insensitive)
	var existingByEmail models.User
	if err := database.DB.Where("LOWER(email) = LOWER(?)", email).First(&existingByEmail).Error; err == nil {
		return nil, errors.New("该邮箱已被注册")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Check for existing user by username (case-insensitive)
	var existingByUsername models.User
	if err := database.DB.Where("LOWER(username) = LOWER(?)", username).First(&existingByUsername).Error; err == nil {
		return nil, errors.New("该用户名已被使用")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	result := database.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("LOWER(email) = LOWER(?)", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func (s *UserService) ValidatePassword(user *models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (s *UserService) UpdateUserPoints(userID uint, points int) error {
	result := database.DB.Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("total_points", gorm.Expr("total_points + ?", points))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
