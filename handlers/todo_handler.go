package handlers

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/Rynoo1/LB-Todo-API/config"
	"github.com/Rynoo1/LB-Todo-API/models"
	"github.com/Rynoo1/LB-Todo-API/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BulkTodoInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

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

// Delete Todo Item
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

// Get Todo Stats for user
func FetchStats(todoRepo *services.TodoService) fiber.Handler {
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

		stats, err := todoRepo.GetUserTodoStats(body.UserId)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to get stats",
			})
		}

		return c.JSON(stats)
	}
}

func parseStatus(s string) models.TodoStatus {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "progress", "in-progress", "in_progress":
		return models.StatusInProgress
	case "done", "completed":
		return models.StatusDone
	default:
		return models.StatusPending
	}
}

// Bulk Upload
func BulkUploadTodos(c *fiber.Ctx) error {
	var body struct {
		UserId uint `json:"user_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "file is required")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to open uploaded file")
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	var inputs []BulkTodoInput

	switch ext {
	case ".json":
		inputs, err = parseJSONTodos(file)
	case ".csv":
		inputs, err = parseCSVTodos(file)
	default:
		return fiber.NewError(fiber.StatusBadRequest, "unsupported file type (use .json or .csv)")
	}

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to parse file: %v", err))
	}

	if len(inputs) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "no todos found in file")
	}

	todos := make([]models.Todo, 0, len(inputs))
	for _, in := range inputs {
		if strings.TrimSpace(in.Title) == "" {
			continue
		}
		todos = append(todos, models.Todo{
			Title:       in.Title,
			Description: in.Description,
			Status:      parseStatus(in.Status),
			UserId:      body.UserId,
		})
	}

	if len(todos) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "no valid todos found in file")
	}

	db := config.DB
	if err := db.Create(&todos).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to save todos")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":       "todos created successfully",
		"created_count": len(todos),
		"user_id":       body.UserId,
	})

}

func parseJSONTodos(r io.Reader) ([]BulkTodoInput, error) {
	var items []BulkTodoInput
	dec := json.NewDecoder(r)
	if err := dec.Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

func parseCSVTodos(r io.Reader) ([]BulkTodoInput, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("empty csv")
	}

	dataRows := rows[1:]

	var items []BulkTodoInput
	for _, row := range dataRows {
		if len(row) == 0 {
			continue
		}

		title := row[0]
		desc := ""
		status := ""

		if len(row) > 1 {
			desc = row[1]
		}
		if len(row) > 2 {
			status = row[2]
		}

		items = append(items, BulkTodoInput{
			Title:       title,
			Description: desc,
			Status:      status,
		})
	}
	return items, nil
}

// All User Todos
func AllUserTodos(todoRepo *services.TodoService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			UserId uint `json:"user_id"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request",
			})
		}

		todos, err := todoRepo.GetUserTodos(body.UserId)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err,
			})
		}

		return c.JSON(todos)
	}
}

// Return Todo
func GetTodo(todoRepo *services.TodoService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			ItemId uint `json:"item_id"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request",
			})
		}

		todo, err := todoRepo.GetTodo(body.ItemId)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err,
			})
		}

		return c.JSON(todo)
	}
}
