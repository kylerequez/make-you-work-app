package main

import (
	"log"
	"time"

	"github.com/kylerequez/make-you-work-app/src/api/db"
	"github.com/kylerequez/make-you-work-app/src/api/models"
	"github.com/kylerequez/make-you-work-app/src/api/repositories"
	"github.com/kylerequez/make-you-work-app/src/api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	log.Println("Seeding database...")
	if err := utils.LoadEnvVariables(); err != nil {
		panic(err)
	}

	if err := db.ConnectDB(); err != nil {
		panic(err)
	}
	conn := db.GetDB("make-you-work-app")

	log.Println("Seeding users...")
	ur := repositories.NewUserRepository(conn, "users")
	if err := SeedUsers(ur); err != nil {
		panic(err)
	}

	// log.Println("Seeding tasks...")
	// tr := repositories.NewTaskRepository(conn, "tasks")
	// seedTask(tr)

	log.Println("Seeding groups...")
	// gr := repositories.NewGroupRepository(conn, "groups")
	// seedGroup(gr)

	log.Println("Seeding complete!")
}

func SeedGroups(gr *repositories.GroupRepository) error {
	return nil
}

func SeedTasks(tr *repositories.TaskRepository) error {
	return nil
}

func SeedUsers(ur *repositories.UserRepository) error {
	password := "password1234"

	hashedPassword, err := utils.HashPassword([]byte(password))
	if err != nil {
		return err
	}

	currentTime := time.Now()

	users := []models.User{
		models.User{
			Firstname:   "Kyle",
			Middlename:  "Pasco",
			Lastname:    "Requez",
			Username:    "keypeearr",
			Email:       "kylerequez155@gmail.com",
			Password:    hashedPassword,
			Status:      utils.USER_STATUS["VERIFIED"],
			Authorities: []string{utils.USER_AUTHORITIES["ADMIN_USER"]},
			CreatedAt:   primitive.NewDateTimeFromTime(currentTime),
			UpdatedAt:   primitive.NewDateTimeFromTime(currentTime),
		},
	}

	for _, user := range users {
		result, err := ur.CreateUser(&user)
		if err != nil {
			return err
		}
		log.Println("Created User with ID: ", result.InsertedID)
	}

	return nil
}
