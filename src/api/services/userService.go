package services

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kylerequez/make-you-work-app/src/api/models"
	"github.com/kylerequez/make-you-work-app/src/api/repositories"
	"github.com/kylerequez/make-you-work-app/src/api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	ur *repositories.UserRepository
}

func NewUserService(ur *repositories.UserRepository) *UserService {
	return &UserService{
		ur: ur,
	}
}

type UserSuccessResponse struct {
	Message string            `json:"message,omitempty"`
	Users   *[]models.UserDTO `json:"users,omitempty"`
	User    *models.UserDTO   `json:"user,omitempty"`
	UserID  any               `json:"_id,omitempty"`
	JWT     *string           `json:"jwt,omitempty"`
}

type UserErrorResponse struct {
	Error string `json:"error"`
}

func (us *UserService) GetAllUsers(c *fiber.Ctx) error {
	results, err := us.ur.GetAllUsers()
	if err != nil {
		return c.Status(http.StatusOK).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(&UserSuccessResponse{
		Users: results,
	})
}

func (us *UserService) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: "id must not be empty",
		})
	}

	oid, err := utils.HexToObjectId(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	result, err := us.ur.GetUserById(oid)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(&UserSuccessResponse{
		User: result,
	})
}

func (us *UserService) CreateUser(c *fiber.Ctx) error {
	type RequestBody struct {
		Firstname   string   `json:"firstname"`
		Middlename  string   `json:"middlename"`
		Lastname    string   `json:"lastname"`
		Email       string   `json:"email"`
		Authorities []string `json:"authorities"`
		Username    string   `json:"username"`
		Password    string   `json:"password"`
	}

	body := new(RequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	// NOTE: TO IMPLEMENT GENERATED PASSWORD THEN SENDING TO EMAIL
	// password, err := utils.GenerateNewPassword()
	// if err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(&UserErrorResponse{
	// 		Error: err.Error(),
	// 	})
	// }

	// hashedPassword, err := utils.HashPassword([]byte(password))
	// if err != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
	// 		Error: err.Error(),
	// 	})
	// }

	hashedPassword, err := utils.HashPassword([]byte(body.Password))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	currentTime := time.Now()
	user := &models.User{
		Firstname:   body.Firstname,
		Middlename:  body.Middlename,
		Lastname:    body.Lastname,
		Username:    body.Username,
		Authorities: body.Authorities,
		Status:      utils.USER_STATUS["NOT_VERIFIED"],
		Email:       body.Email,
		Password:    hashedPassword,
		CreatedAt:   primitive.NewDateTimeFromTime(currentTime),
		UpdatedAt:   primitive.NewDateTimeFromTime(currentTime),
	}

	result, err := us.ur.CreateUser(user)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(&UserSuccessResponse{
		UserID: result.InsertedID,
	})
}

func (us *UserService) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: "id must not be empty",
		})
	}

	oid, err := utils.HexToObjectId(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	type RequestBody struct {
		Firstname       string   `json:"firstname"`
		Middlename      string   `json:"middlename"`
		Lastname        string   `json:"lastname"`
		Email           string   `json:"email"`
		Authorities     []string `json:"authorities"`
		Status          string   `json:"status"`
		Username        string   `json:"username"`
		Password        string   `json:"password"`
		ConfirmPassword string   `json:"confirmPassword"`
	}

	body := new(RequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	if err := utils.ComparePassword(body.Password, body.ConfirmPassword); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	hashedPassword, err := utils.HashPassword([]byte(body.Password))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	updatedUser := &models.User{
		Firstname:   body.Firstname,
		Middlename:  body.Middlename,
		Lastname:    body.Lastname,
		Authorities: body.Authorities,
		Status:      body.Status,
		Username:    body.Username,
		Email:       body.Email,
		Password:    hashedPassword,
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	result, err := us.ur.UpdateUser(oid, updatedUser)
	if err != nil || result == nil {
		return c.Status(http.StatusInternalServerError).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(&UserSuccessResponse{
		UserID: result.UpsertedID,
	})
}

func (us *UserService) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: "id must not be empty",
		})
	}

	oid, err := utils.HexToObjectId(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	_, err = us.ur.DeleteUser(oid)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(nil)
}

func (us *UserService) LoginUser(c *fiber.Ctx) error {
	type RequestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	body := new(RequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	result, err := us.ur.GetUserByEmailOrUsername(body.Username, body.Username)
	if err != nil && result == nil {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	if result.Status == utils.USER_STATUS["NOT_VERIFIED"] || result.Status == utils.USER_STATUS["DEACTIVATED"] {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: "please check your email and verify your account to log in",
		})
	}

	if err := utils.ValidatePassword([]byte(body.Password), result.Password); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	jwtString, err := utils.GenerateJWT(result)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	jwtCookie := fiber.Cookie{
		Name:     "make-you-work-auth-token",
		Value:    *jwtString,
		Expires:  time.Now().Add(15 * time.Minute),
		Secure:   true,
		HTTPOnly: true,
	}
	c.Cookie(&jwtCookie)

	return c.Status(http.StatusOK).JSON(&UserSuccessResponse{
		JWT: jwtString,
	})
}

func (us *UserService) RegisterUser(c *fiber.Ctx) error {
	type RequestBody struct {
		Firstname       string `json:"firstname"`
		Middlename      string `json:"middlename"`
		Lastname        string `json:"lastname"`
		Email           string `json:"email"`
		Status          string `json:"status"`
		Username        string `json:"username"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword" `
	}

	body := new(RequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	if err := utils.ComparePassword(body.Password, body.ConfirmPassword); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	hashedPassword, err := utils.HashPassword([]byte(body.Password))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	user := &models.User{
		Firstname:   body.Firstname,
		Middlename:  body.Middlename,
		Lastname:    body.Lastname,
		Username:    body.Username,
		Authorities: []string{utils.USER_AUTHORITIES["NORMAL_USER"]},
		Status:      utils.USER_STATUS["NOT_VERIFIED"],
		Email:       body.Email,
		Password:    hashedPassword,
	}

	result, err := us.ur.CreateUser(user)
	if err != nil || result == nil {
		return c.Status(http.StatusBadRequest).JSON(&UserErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(&UserSuccessResponse{
		Message: "successfully registered a new user",
	})
}
