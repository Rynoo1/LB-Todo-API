package routes

import (
	"github.com/Rynoo1/LB-Todo-API/handlers"
	"github.com/Rynoo1/LB-Todo-API/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, svc *services.AppServices, db *gorm.DB, authService *services.AuthService) {

	authHandler := handlers.NewAuthHandler(db, authService, svc.UserServices)

	app.Post("/auth/register", authHandler.Register)
	app.Post("/auth/login", authHandler.Login)

	// Protected routes

	// Creat Todo  																				V
	app.Post("/todo/create", func(c *fiber.Ctx) error {
		return handlers.CreateTodo(c, svc)
	})

	/* TODOS */
	// Update Todo Status  																		V
	app.Post("/todo/update/status", handlers.UpdateTodoStatus(svc.TodoServices))
	// Update Todo Description  																V
	app.Post("/todo/update/desc", handlers.UpdateTodoDesc(svc.TodoServices))
	// Update Todo Title  																		V
	app.Post("/todo/update/title", handlers.UpdateTodoTitle(svc.TodoServices))
	// Delete  																					TODO
	app.Post("/todo/delete", handlers.DeleteTodoItem(svc.TodoServices))

	/* USER */
	// Soft delete account  																	TODO
	app.Post("/user/delete", handlers.SoftDeleteUser(svc.UserServices))
	// Soft delete account and todos  															TODO
	app.Post("/user/deleteall", handlers.SoftDeleteUserWithTodos(svc.UserServices))
}
