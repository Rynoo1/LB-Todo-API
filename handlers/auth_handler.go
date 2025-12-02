package handlers

import (
	"net/http"

	"github.com/Rynoo1/LB-Todo-API/models"
	"github.com/Rynoo1/LB-Todo-API/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db           *gorm.DB
	authServices *services.AuthService
	userService  *services.UserService
}

func NewAuthHandler(db *gorm.DB, authService *services.AuthService, userService *services.UserService) *AuthHandler {
	return &AuthHandler{
		db:           db,
		authServices: authService,
		userService:  userService,
	}
}

type RegisterRrequest struct {
	Name     string `json:"name" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password" validate:"required, min=6"`
	Username string `json:"username" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Message string       `json:"message"`
	Token   string       `json:"token"`
	User    *models.User `json:"user"`
}

// Register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRrequest
	err := c.BodyParser(&req)
	if err != nil {
		c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
		return err
	}

	user, err := h.userService.CreateUser(req.Name, req.Surname, req.Username, req.Email, req.Password)
	if err != nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}

	// Generate token
	token, err := h.authServices.GenerateToken(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.Status(http.StatusCreated).JSON(AuthResponse{
		Message: "User Registered successfully",
		Token:   token,
		User:    user,
	})
}

// Login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	err := c.BodyParser(&req)
	if err != nil {
		c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
		return err
	}

	user, err := h.userService.FindByEmail(req.Email)
	if err != nil || !user.CheckPassword(req.Password) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token, err := h.authServices.GenerateToken(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(AuthResponse{
		Message: "Login successful",
		Token:   token,
		User:    user,
	})

}
