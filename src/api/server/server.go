package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kylerequez/make-you-work-app/src/api/controllers"
	"github.com/kylerequez/make-you-work-app/src/api/utils"
)

func InitServer() {
	if err := utils.LoadEnvVariables(); err != nil {
		panic(err)
	}

	// engine := html.New("./src/ui/views", ".html")
	// engine.Reload(true)

	server := fiber.New(
		fiber.Config{
			AppName: "Make You Work",
			// Views:         engine,
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

	if err := controllers.InitRoutes(server); err != nil {
		panic(err)
	}

	server.Get("/", func(c *fiber.Ctx) error {
		// main := layouts.MainLayout(" | Home", "")
		// handler := adaptor.HTTPHandler(templ.Handler(main))
		// return handler(c)

		err := utils.SendEmail(
			"Make You Work App Registration",
			"kyle.requez.test@gmail.com",
			[]string{"kylerequez155@gmail.com"},
			nil,
			"registration.html",
			map[string]string{
				"Title": "Title",
			},
		)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"message": "An email was sent to your inbox. Please check it.",
		})
	})

	serverPort, err := utils.GetEnv("SERVER_PORT")
	if err != nil {
		panic("::: Error in loading the server port: " + err.Error())
	}
	server.Listen(":" + *serverPort)
}
