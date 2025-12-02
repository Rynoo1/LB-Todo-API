package services

import (
	"fmt"

	"github.com/Rynoo1/LB-Todo-API/models"
	"gorm.io/gorm"
)

type TodoService struct {
	db *gorm.DB
}

func NewTodoService(db *gorm.DB) *TodoService {
	svc := &TodoService{
		db: db,
	}
	return svc
}

// Create Todo
func (s *TodoService) CreateTodo(todo *models.Todo) error {
	return s.db.Create(todo).Error
}

// Update Todo Status
func (s *TodoService) UpdateStatus(itemId, userId uint, status models.TodoStatus) error {
	if !status.IsValid() {
		return fmt.Errorf("invalid status")
	}

	result := s.db.Model(&models.Todo{}).
		Where("id = ? AND user_id = ?", itemId, userId).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Update Todo description
func (s *TodoService) UpdateDesc(itemId, userId uint, desc string) error {

	result := s.db.Model(&models.Todo{}).Where("id = ? AND user_id = ?", itemId, userId).Update("description", desc)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Update Todo Title
func (s *TodoService) UpdateTitle(itemId, userId uint, title string) error {
	result := s.db.Model(&models.Todo{}).Where("id = ? AND user_id = ?", itemId, userId).Update("title", title)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s *TodoService) DeleteTodo(itemId, userId uint) error {
	result := s.db.Where("id = ? AND user_id = ?", itemId, userId).Delete(&models.Todo{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete Todo item
// Upload JSON/CSV file
