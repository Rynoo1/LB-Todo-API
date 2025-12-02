package services

import (
	"errors"

	"github.com/Rynoo1/LB-Todo-API/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(name, surname, username, email, password string) (*models.User, error) {
	var existing models.User

	err := s.db.Where("email = ?", email).First(&existing).Error
	if err == nil {
		return nil, errors.New("user already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user := models.User{
		Name:     name,
		Surname:  surname,
		Username: username,
		Email:    email,
		Password: password,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Find User by Email
func (s *UserService) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) SoftDeleteUser(userId uint) error {
	result := s.db.Delete(&models.User{}, userId)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s *UserService) SoftDeleteUserWithTodos(userId uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.User{}, userId).Error; err != nil {
			return err
		}

		if err := tx.Where("user_id = ?", userId).Delete(&models.Todo{}).Error; err != nil {
			return err
		}
		return nil
	})
}
