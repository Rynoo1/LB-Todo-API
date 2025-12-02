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

type TodoStats struct {
	UserId     uint  `json:"user_id"`
	Total      int64 `json:"total"`
	Pending    int64 `json:"pending"`
	InProgress int64 `json:"progress"`
	Done       int64 `json:"done"`
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

func (s *TodoService) GetUserTodoStats(userId uint) (*TodoStats, error) {
	var total int64
	var pending int64
	var inProgress int64
	var done int64

	if err := s.db.Model(&models.Todo{}).Where("user_id = ?", userId).Count(&total).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.Todo{}).Where("user_id = ? AND status = ?", userId, models.StatusPending).Count(&pending).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.Todo{}).Where("user_id = ? AND status = ?", userId, models.StatusInProgress).Count(&inProgress).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.Todo{}).Where("user_id = ? AND status = ?", userId, models.StatusDone).Count(&done).Error; err != nil {
		return nil, err
	}

	return &TodoStats{
		UserId:     userId,
		Total:      total,
		Pending:    pending,
		InProgress: inProgress,
		Done:       done,
	}, nil
}

// Todos Analytics
// Upload JSON/CSV file
