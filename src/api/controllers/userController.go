package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kylerequez/make-you-work-app/src/api/middlewares"
	"github.com/kylerequez/make-you-work-app/src/api/services"
)

type UserController struct {
	us *services.UserService
}

func NewUserController(us *services.UserService) *UserController {
	return &UserController{
		us: us,
	}
}

func (uc *UserController) InitUserController(server *fiber.App) {
	userRoutes := server.Group("/api/v1")
	userRoutes.Get("/users", uc.GetAllUsers)
	userRoutes.Get("/users/:id", uc.GetUserById)

	protected := userRoutes.Group("/users",
		middlewares.AuthenticateJwtToken,
		middlewares.CheckIfAdminUser,
	)
	protected.Post("", uc.CreateUser)
	protected.Patch("/:id", uc.UpdateUser)
	protected.Delete("/:id", uc.DeleteUser)
}

func (uc *UserController) GetAllUsers(c *fiber.Ctx) error {
	return uc.us.GetAllUsers(c)
}

func (uc *UserController) GetUserById(c *fiber.Ctx) error {
	return uc.us.GetUserById(c)
}

func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	return uc.us.CreateUser(c)
}

func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	return uc.us.UpdateUser(c)
}

func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	return uc.us.DeleteUser(c)
}
