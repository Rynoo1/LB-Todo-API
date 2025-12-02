package routes

import (
	"github.com/Rynoo1/LB-Todo-API/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, svc *services.AppServices, db *gorm.DB, authService *services.AuthService) {

	// authHandler

}
