package server

import (
	"os"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/kylerequez/make-you-work-app/src/api/controllers"
	"github.com/kylerequez/make-you-work-app/src/api/utils"
	"github.com/kylerequez/make-you-work-app/src/ui/views/layouts"
)

func InitServer() {
	if err := utils.LoadEnvVariables(); err != nil {
		panic(err)
	}

	engine := html.New("./src/ui/views", ".html")
	engine.Reload(true)

	server := fiber.New(
		fiber.Config{
			AppName:       "Make You Work",
			Views:         engine,
			CaseSensitive: true,
		},
	)

	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,DELETE,PATCH",
	}))
	server.Use(logger.New(logger.Config{
		Format:   "[${ip}]:${port} ${status} - ${method} ${path}\n",
		TimeZone: "Asia/Manila",
	}))

	controllers.InitRoutes(server)

	server.Get("/", func(c *fiber.Ctx) error {
		main := layouts.MainLayout(" | Home", "")
		handler := adaptor.HTTPHandler(templ.Handler(main))
		return handler(c)

	})

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		panic("::: Error in loading the server port")
	}
	server.Listen(":" + serverPort)
}
