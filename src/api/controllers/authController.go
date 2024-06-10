package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kylerequez/make-you-work-app/src/api/services"
)

type AuthController struct {
	us *services.UserService
}

func NewAuthController(us *services.UserService) *AuthController {
	return &AuthController{
		us: us,
	}
}

func (ac *AuthController) InitAuthController(server *fiber.App) {
	authRoutes := server.Group("/api/v1/auth")
	authRoutes.Post("/login", ac.LoginUser)
	authRoutes.Post("/register", ac.RegisterUser)
}

func (ac *AuthController) LoginUser(c *fiber.Ctx) error {
	return ac.us.LoginUser(c)
}

func (ac *AuthController) RegisterUser(c *fiber.Ctx) error {
	return ac.us.RegisterUser(c)
}
