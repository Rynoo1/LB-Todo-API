package services

import (
	"github.com/Rynoo1/LB-Todo-API/models"
	"gorm.io/gorm"
)

type TodoService struct {
	db *gorm.DB
}

func (s *TodoService) CreateTodo(todo *models.Todo) error {
	return s.db.Create(todo).Error
}

// Update Todo Status
// Change Description
// Change Title
// Delete Todo item
// Upload JSON/CSV file
