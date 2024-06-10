package services

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kylerequez/make-you-work-app/src/api/models"
	"github.com/kylerequez/make-you-work-app/src/api/repositories"
	"github.com/kylerequez/make-you-work-app/src/api/utils"
)

type GroupService struct {
	gr *repositories.GroupRepository
}

func NewGroupService(gr *repositories.GroupRepository) *GroupService {
	return &GroupService{
		gr: gr,
	}
}

type GroupSuccessResponse struct {
	Message string             `json:"message,omitempty"`
	Groups  *[]models.GroupDTO `json:"groups,omitempty"`
	Group   *models.GroupDTO   `json:"group,omitempty"`
	GroupID any                `json:"_id,omitempty"`
}

type GroupErrorResponse struct {
	Error string `json:"error"`
}

func (gs *GroupService) GetAllGroups(c *fiber.Ctx) error {
	results, err := gs.gr.GetAllGroups()
	if err != nil {
		return c.Status(http.StatusOK).JSON(&GroupErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(&GroupSuccessResponse{
		Groups: results,
	})
}

func (gs *GroupService) GetGroupById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(&GroupErrorResponse{
			Error: "id must not be empty",
		})
	}

	oid, err := utils.HexToObjectId(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&GroupErrorResponse{
			Error: err.Error(),
		})
	}

	result, err := gs.gr.GetGroupById(oid)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&GroupErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(&GroupSuccessResponse{
		Group: result,
	})
}
