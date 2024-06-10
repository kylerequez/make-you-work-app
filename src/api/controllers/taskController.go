package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kylerequez/make-you-work-app/src/api/services"
)

type TaskController struct {
	ts *services.TaskService
}

func NewTaskController(ts *services.TaskService) *TaskController {
	return &TaskController{
		ts: ts,
	}
}

func (tc *TaskController) InitTaskController(server *fiber.App) {
	taskRoutes := server.Group("/api/v1")
	taskRoutes.Get("/tasks", tc.GetAllTasks)
	taskRoutes.Get("/tasks/:id", tc.GetTaskById)
	taskRoutes.Get("/tasks/user/:id", tc.GetTasksByUserId)
	taskRoutes.Post("/tasks", tc.CreateTask)
	taskRoutes.Patch("/tasks/:id", tc.UpdateTask)
	taskRoutes.Delete("/tasks/:id", tc.DeleteTask)
}

func (tc *TaskController) GetAllTasks(c *fiber.Ctx) error {
	return tc.ts.GetAllTasks(c)
}

func (tc *TaskController) GetTaskById(c *fiber.Ctx) error {
	return tc.ts.GetTaskById(c)
}

func (tc *TaskController) GetTasksByUserId(c *fiber.Ctx) error {
	return tc.ts.GetTasksByUserId(c)
}

func (tc *TaskController) CreateTask(c *fiber.Ctx) error {
	return tc.ts.CreateTask(c)
}

func (tc *TaskController) UpdateTask(c *fiber.Ctx) error {
	return tc.ts.UpdateTask(c)
}

func (tc *TaskController) DeleteTask(c *fiber.Ctx) error {
	return tc.ts.DeleteTask(c)
}
