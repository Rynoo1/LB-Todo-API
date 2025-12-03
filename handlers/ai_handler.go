package handlers

import (
	"encoding/json"

	"github.com/Rynoo1/LB-Todo-API/services"
	"github.com/gofiber/fiber/v2"
)

func StoryPoints(repo *services.AppServices) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			ItemId uint `json:"item_id"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request",
			})
		}

		todoItem, err := repo.TodoServices.GetTodo(body.ItemId)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err,
			})
		}

		aiResponse, err := repo.AiServices.AiServiceCall(todoItem.Title, todoItem.Description, "storypoints")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err,
			})
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(aiResponse), &parsed); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "AI returned invalid JSON",
			})
		}

		return c.JSON(parsed)
	}
}

func TodoSteps(repo *services.AppServices) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			ItemId uint `json:"item_id"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request",
			})
		}

		todoItem, err := repo.TodoServices.GetTodo(body.ItemId)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err,
			})
		}

		aiResponse, err := repo.AiServices.AiServiceCall(todoItem.Title, todoItem.Description, "todosteps")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err,
			})
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(aiResponse), &parsed); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "AI returned invalid JSON",
			})
		}

		return c.JSON(parsed)
	}
}
