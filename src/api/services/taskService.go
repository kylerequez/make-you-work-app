package services

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/kylerequez/make-you-work-app/src/api/models"
	"github.com/kylerequez/make-you-work-app/src/api/repositories"
	"github.com/kylerequez/make-you-work-app/src/api/utils"
	"github.com/kylerequez/make-you-work-app/src/ui/views"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService struct {
	tr *repositories.TaskRepository
}

func NewTaskService(tr *repositories.TaskRepository) *TaskService {
	return &TaskService{
		tr: tr,
	}
}

type TaskSuccessResponse struct {
	Message string            `json:"message,omitempty"`
	Tasks   *[]models.TaskDTO `json:"tasks,omitempty"`
	Task    *models.TaskDTO   `json:"task,omitempty"`
	TaskID  any               `json:"_id,omitempty"`
}

type TaskErrorResponse struct {
	Error string `json:"error"`
}

func (ts *TaskService) ViewGetAllTasks(c *fiber.Ctx) error {
	results, err := ts.tr.GetAllTasks()
	if err != nil {
		tasks := views.Tasks("| View Tasks", nil, err.Error())
		handler := adaptor.HTTPHandler(templ.Handler(tasks))
		return handler(c)
	}

	tasks := views.Tasks("| View Tasks", *results, "")
	handler := adaptor.HTTPHandler(templ.Handler(tasks))
	return handler(c)
}

func (ts *TaskService) GetAllTasks(c *fiber.Ctx) error {
	results, err := ts.tr.GetAllTasks()
	if err != nil {
		return c.Status(http.StatusOK).JSON(&TaskErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(&TaskSuccessResponse{
		Tasks: results,
	})
}

func (ts *TaskService) GetTaskById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(&TaskErrorResponse{
			Error: "id must not be empty",
		})
	}

	oid, err := utils.HexToObjectId(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&TaskErrorResponse{
			Error: err.Error(),
		})
	}

	result, err := ts.tr.GetTaskById(oid)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&TaskErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(&TaskSuccessResponse{
		Task: result,
	})
}

func (ts *TaskService) GetTasksByUserId(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(&TaskErrorResponse{
			Error: "id must not be empty",
		})
	}

	oid, err := utils.HexToObjectId(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&TaskErrorResponse{
			Error: err.Error(),
		})
	}

	result, err := ts.tr.GetTasksByUserId(oid)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&TaskErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(&TaskSuccessResponse{
		Tasks: result,
	})
}

func (ts *TaskService) CreateTask(c *fiber.Ctx) error {
	type RequestBody struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	body := new(RequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&TaskErrorResponse{
			Error: err.Error(),
		})
	}

	jwt := c.Cookies("make-you-work-auth-token")
	if jwt == "" {
		return c.Status(http.StatusUnauthorized).JSON(&TaskErrorResponse{
			Error: "the auth cookie does not exist",
		})
	}

	id, err := utils.RetrieveUserIdFromJwt(jwt)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&TaskErrorResponse{
			Error: err.Error(),
		})
	}

	currentTime := time.Now()
	task := &models.Task{
		Title:       body.Title,
		Description: body.Description,
		Status:      utils.TASK_STATUS["ON-GOING"],
		CreatedBy:   id,
		CreatedAt:   primitive.NewDateTimeFromTime(currentTime),
		UpdatedAt:   primitive.NewDateTimeFromTime(currentTime),
		CompletedAt: nil,
	}

	result, err := ts.tr.CreateTask(task)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&TaskErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(&TaskSuccessResponse{
		TaskID: result.InsertedID,
	})
}

func (ts *TaskService) UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(&TaskErrorResponse{
			Error: "id must not be empty",
		})
	}

	oid, err := utils.HexToObjectId(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&TaskErrorResponse{
			Error: err.Error(),
		})
	}

	type RequestBody struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	body := new(RequestBody)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&TaskErrorResponse{
			Error: err.Error(),
		})
	}

	currentTime := primitive.NewDateTimeFromTime(time.Now())
	var completedAt *primitive.DateTime = nil
	if body.Status == utils.TASK_STATUS["COMPLETED"] {
		completedAt = &currentTime
	}

	updatedTask := &models.Task{
		Title:       body.Title,
		Description: body.Description,
		Status:      body.Status,
		CompletedAt: completedAt,
		UpdatedAt:   currentTime,
	}

	result, err := ts.tr.UpdateTask(oid, updatedTask)
	if err != nil || result == nil {
		return c.Status(http.StatusInternalServerError).JSON(&TaskErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(&TaskSuccessResponse{
		TaskID: result.UpsertedID,
	})
}

func (ts *TaskService) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(&TaskErrorResponse{
			Error: "id must not be empty",
		})
	}

	oid, err := utils.HexToObjectId(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&TaskErrorResponse{
			Error: err.Error(),
		})
	}

	_, err = ts.tr.DeleteTask(oid)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&TaskErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(nil)
}
