package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kylerequez/make-you-work-app/src/api/middlewares"
	"github.com/kylerequez/make-you-work-app/src/api/services"
)

type GroupController struct {
	gs *services.GroupService
}

func NewGroupController(gs *services.GroupService) *GroupController {
	return &GroupController{
		gs: gs,
	}
}

func (gc *GroupController) InitGroupController(server *fiber.App) {
	GroupRoutes := server.Group("/api/v1", middlewares.AuthenticateJwtToken)
	GroupRoutes.Get("/groups", gc.GetAllGroups)
	GroupRoutes.Get("/groups/:id", gc.GetGroupById)
	// GroupRoutes.Get("/groups/user/:id", tc.GetGroupsByUserId)
	// GroupRoutes.Post("/groups", tc.CreateGroup)
	// GroupRoutes.Patch("/groups/:id", tc.UpdateGroup)
	// GroupRoutes.Delete("/groups/:id", tc.DeleteGroup)
}

func (gc *GroupController) GetAllGroups(c *fiber.Ctx) error {
	return gc.gs.GetAllGroups(c)
}

func (gc *GroupController) GetGroupById(c *fiber.Ctx) error {
	return gc.gs.GetGroupById(c)
}
