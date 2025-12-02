package handlers

import (
	"strings"
	"time"

	"github.com/Rynoo1/LB-Todo-API/models"
	"github.com/Rynoo1/LB-Todo-API/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateTodo(c *fiber.Ctx, todoRepo *services.AppServices) error {
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		UserId      uint   `json:"user_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	todo := models.Todo{Title: body.Title, Description: body.Description, UserId: body.UserId, Created_at: time.Now(), Updated_at: time.Now(), Status: "pending"}

	if err := todoRepo.TodoServices.CreateTodo(&todo); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to create todo item",
		})
	}

	return c.Status(201).JSON(todo)
}

func UpdateTodoStatus(todoRepo *services.TodoService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			ItemId uint   `json:"item_id"`
			UserId uint   `json:"user_id"`
			Status string `json:"status"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request",
			})
		}

		status, err := models.ParseTodoStatus(body.Status)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = todoRepo.UpdateStatus(body.ItemId, body.UserId, status)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(404).JSON(fiber.Map{
					"error": "todo not found",
				})
			}
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to update status",
			})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}

}

// Update Todo Description
func UpdateTodoDesc(todoRepo *services.TodoService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			ItemId uint    `json:"item_id"`
			UserId uint    `json:"user_id"`
			Desc   *string `json:"description"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request",
			})
		}

		if body.Desc == nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "description is required",
			})
		}

		if body.ItemId == 0 {
			return c.Status(400).JSON(fiber.Map{
				"error": "item_id is required",
			})
		}

		err := todoRepo.UpdateDesc(body.ItemId, body.UserId, *body.Desc)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(404).JSON(fiber.Map{
					"error": "todo not found",
				})
			}
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to update description",
			})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// Update Todo Title
func UpdateTodoTitle(todoRepo *services.TodoService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			ItemId uint   `json:"item_id"`
			UserId uint   `json:"user_id"`
			Title  string `json:"title"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request",
			})
		}

		if strings.TrimSpace(body.Title) == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "title is required",
			})
		}

		if body.ItemId == 0 {
			return c.Status(400).JSON(fiber.Map{
				"error": "item_id is required",
			})
		}

		err := todoRepo.UpdateTitle(body.ItemId, body.UserId, body.Title)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(400).JSON(fiber.Map{
					"error": "todo not found",
				})
			}
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to update title",
			})
		}
		return c.Status(201).JSON(fiber.Map{
			"success": "title updated",
		})
	}

}

func DeleteTodoItem(todoRepo *services.TodoService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			UserId uint `json:"user_id"`
			ItemId uint `json:"item_id"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request",
			})
		}

		if body.ItemId == 0 {
			return c.Status(400).JSON(fiber.Map{
				"error": "item_id is required",
			})
		}

		if body.UserId == 0 {
			return c.Status(400).JSON(fiber.Map{
				"error": "user_id is required",
			})
		}

		err := todoRepo.DeleteTodo(body.ItemId, body.UserId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(404).JSON(fiber.Map{
					"error": "todo item not found",
				})
			}
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to delete todo item",
			})
		}

		return c.Status(201).JSON(fiber.Map{
			"success": "item successfully deleted",
		})
	}
}
