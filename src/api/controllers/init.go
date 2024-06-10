package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kylerequez/make-you-work-app/src/api/db"
	"github.com/kylerequez/make-you-work-app/src/api/repositories"
	"github.com/kylerequez/make-you-work-app/src/api/services"
)

func InitRoutes(server *fiber.App) error {
	if err := db.ConnectDB(); err != nil {
		return err
	}
	conn := db.GetDB("make-you-work-app")

	ur := repositories.NewUserRepository(conn, "users")
	us := services.NewUserService(ur)

	ac := NewAuthController(us)
	ac.InitAuthController(server)

	uc := NewUserController(us)
	uc.InitUserController(server)

	tr := repositories.NewTaskRepository(conn, "tasks")
	ts := services.NewTaskService(tr)
	tc := NewTaskController(ts)
	tc.InitTaskController(server)

	gr := repositories.NewGroupRepository(conn, "groups")
	gs := services.NewGroupService(gr)
	gc := NewGroupController(gs)
	gc.InitGroupController(server)

	return nil
}
