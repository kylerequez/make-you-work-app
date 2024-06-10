package middlewares

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kylerequez/make-you-work-app/src/api/utils"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func AuthenticateJwtToken(c *fiber.Ctx) error {
	jwt := c.Cookies("make-you-work-auth-token")
	if jwt == "" {
		return c.Status(http.StatusUnauthorized).JSON(&ErrorResponse{
			Error: "the auth cookie does not exist",
		})
	}

	isValid, err := utils.ValidateJwt(jwt)
	if !isValid {
		return c.Status(http.StatusUnauthorized).JSON(&ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Next()
}

func CheckIfAdminUser(c *fiber.Ctx) error {
	jwt := c.Cookies("make-you-work-auth-token")
	if jwt == "" {
		return c.Status(http.StatusUnauthorized).JSON(&ErrorResponse{
			Error: "the auth cookie does not exist",
		})
	}

	authorities := [1]string{utils.USER_AUTHORITIES["ADMIN_USER"]}
	hasAuth, err := utils.CheckIfUserHasAuthorities(jwt, authorities[:])
	if !hasAuth {
		return c.Status(http.StatusUnauthorized).JSON(&ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Next()
}

func CheckIfNormalUser(c *fiber.Ctx) error {
	jwt := c.Cookies("make-you-work-auth-token")
	if jwt == "" {
		return c.Status(http.StatusUnauthorized).JSON(&ErrorResponse{
			Error: "the auth cookie does not exist",
		})
	}

	authorities := [1]string{utils.USER_AUTHORITIES["NORMAL_USER"]}
	hasAuth, err := utils.CheckIfUserHasAuthorities(jwt, authorities[:])
	if !hasAuth {
		return c.Status(http.StatusUnauthorized).JSON(&ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Next()
}
