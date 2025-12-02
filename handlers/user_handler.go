package handlers

import (
	"github.com/Rynoo1/LB-Todo-API/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SoftDeleteUser(userRepo *services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			UserId uint `json:"user_id"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request",
			})
		}

		if err := userRepo.SoftDeleteUser(body.UserId); err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "user not found",
				})
			}
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to delete user",
			})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func SoftDeleteUserWithTodos(userRepo *services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			UserId uint `json:"user_id"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request",
			})
		}
		if body.UserId == 0 {
			return c.Status(400).JSON(fiber.Map{
				"error": "user_id is required",
			})
		}

		err := userRepo.SoftDeleteUserWithTodos(body.UserId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(404).JSON(fiber.Map{
					"error": "user not found",
				})
			}
			return c.Status(500).JSON(fiber.Map{
				"error": "user and todos could not be deleted",
			})
		}

		return c.Status(201).JSON(fiber.Map{
			"success": "user and todos deleted",
		})
	}
}
