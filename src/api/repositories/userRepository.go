package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/kylerequez/make-you-work-app/src/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewUserRepository(db *mongo.Database, collectionName string) *UserRepository {
	return &UserRepository{
		db:             db,
		collectionName: collectionName,
	}
}

func (ur *UserRepository) GetAllUsers() (*[]models.UserDTO, error) {
	ctx := context.TODO()
	filter := bson.D{{}}
	projection := bson.D{{"password", 0}}
	opts := options.Find().SetProjection(projection)

	cursor, err := ur.db.Collection(ur.collectionName).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.UserDTO
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	return &results, nil
}

func (ur *UserRepository) GetUserById(id *primitive.ObjectID) (*models.UserDTO, error) {
	ctx := context.TODO()
	filter := bson.D{{"_id", id}}

	var result models.UserDTO
	err := ur.db.Collection(ur.collectionName).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *UserRepository) CreateUser(newUser *models.User) (*mongo.InsertOneResult, error) {
	ctx := context.TODO()

	isExists, _ := ur.GetUserByEmailOrUsername(newUser.Email, newUser.Username)
	if isExists != nil {
		return nil, errors.New("the email/username exists")
	}

	result, err := ur.db.Collection(ur.collectionName).InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ur *UserRepository) GetUserByEmailOrUsername(email string, username string) (*models.User, error) {
	ctx := context.TODO()
	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{"email", email}},
				bson.D{{"username", username}},
			},
		},
	}

	var result models.User
	err := ur.db.Collection(ur.collectionName).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *UserRepository) UpdateUser(id *primitive.ObjectID, user *models.User) (*mongo.UpdateResult, error) {
	isExists, err := ur.GetUserById(id)
	if isExists == nil {
		return nil, err
	}

	isExistsCredentials, _ := ur.GetUserByEmailOrUsername(user.Email, user.Username)
	if isExistsCredentials != nil {
		return nil, errors.New("the email/username exists")
	}

	ctx := context.TODO()
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"firstname", user.Firstname},
			{"middlename", user.Middlename},
			{"lastname", user.Lastname},
			{"username", user.Username},
			{"email", user.Email},
			{"password", user.Password},
			{"updatedAt", primitive.NewDateTimeFromTime(time.Now())},
		},
		},
	}

	result, err := ur.db.Collection(ur.collectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ur *UserRepository) DeleteUser(id *primitive.ObjectID) (*mongo.DeleteResult, error) {
	ctx := context.TODO()
	filter := bson.D{{"_id", id}}

	result, err := ur.db.Collection(ur.collectionName).DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}
