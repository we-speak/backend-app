package postgres

import (
	"backend-app/internal/config"
	"backend-app/internal/storage/models"

	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	DB *gorm.DB
}

func New(cfg *config.Config) (Storage, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})

	if err != nil {
		return Storage{}, err
	}
	db.AutoMigrate(models.User{})
	return Storage{DB: db}, nil
}

func (s *Storage) CreateUser(user *models.User) error {
	if err := s.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Storage) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Storage) UpdateUser(user *models.User) error {
	if err := s.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteUser(id uint) error {
	if err := s.DB.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetAllUsers(offset int, limit int) ([]models.User, error) {
	var users []models.User
	if err := s.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
