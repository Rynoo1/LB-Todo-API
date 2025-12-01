package migrate

import (
	"fmt"

	"github.com/Rynoo1/LB-Todo-API/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Todo{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}
	return nil
}
