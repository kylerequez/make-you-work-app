package repositories

import (
	"context"
	"log"
	"time"

	"github.com/kylerequez/make-you-work-app/src/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewTaskRepository(db *mongo.Database, collectionName string) *TaskRepository {
	return &TaskRepository{
		db:             db,
		collectionName: collectionName,
	}
}

func (tr *TaskRepository) GetAllTasks() (*[]models.TaskDTO, error) {
	ctx := context.TODO()

	lookupStage := bson.D{
		{"$lookup", bson.D{
			{
				"from",
				"users",
			},
			{
				"localField",
				"createdBy",
			},
			{
				"foreignField",
				"_id",
			},
			{
				"as",
				"createdBy",
			},
		},
		}}

	pipeline := mongo.Pipeline{lookupStage}

	cursor, err := tr.db.Collection(tr.collectionName).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.TaskDTO
	for cursor.Next(ctx) {
		var taskWithUser struct {
			models.TaskDTO `bson:"inline"`
			CreatedBy      []models.UserDTO `bson:"createdBy"`
		}
		err := cursor.Decode(&taskWithUser)
		if err != nil {
			return nil, err
		}

		if len(taskWithUser.CreatedBy) > 0 {
			taskWithUser.TaskDTO.CreatedBy = taskWithUser.CreatedBy[0]
		}
		results = append(results, taskWithUser.TaskDTO)
	}

	return &results, nil
}

func (tr *TaskRepository) GetTaskById(id *primitive.ObjectID) (*models.TaskDTO, error) {
	ctx := context.TODO()

	matchStage := bson.D{
		{"$match", bson.D{
			{
				"_id", id,
			},
		}},
	}

	lookupCreatedAtStage := bson.D{
		{"$lookup", bson.D{
			{
				"from",
				"users",
			},
			{
				"localField",
				"createdBy",
			},
			{
				"foreignField",
				"_id",
			},
			{
				"as",
				"createdBy",
			},
		}},
	}

	// lookupAssignedForStage := bson.D{
	// 	{"$lookup", bson.D{
	// 		{
	// 			"from",
	// 			"users",
	// 		},
	// 		{
	// 			"localField",
	// 			"assignedFor",
	// 		},
	// 		{
	// 			"foreignField",
	// 			"_id",
	// 		},
	// 		{
	// 			"as",
	// 			"assignedFor",
	// 		},
	// 	}},
	// }

	pipeline := mongo.Pipeline{matchStage, lookupCreatedAtStage}

	cursor, err := tr.db.Collection(tr.collectionName).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result models.TaskDTO
	for cursor.Next(ctx) {
		var taskWithUser struct {
			models.TaskDTO `bson:"inline"`
			CreatedBy      []models.UserDTO `bson:"createdBy"`
			// AssignedFor
		}

		log.Println(taskWithUser)
		err := cursor.Decode(&taskWithUser)
		if err != nil {
			return nil, err
		}

		if len(taskWithUser.CreatedBy) > 0 {
			taskWithUser.TaskDTO.CreatedBy = taskWithUser.CreatedBy[0]
		}

		result = taskWithUser.TaskDTO
	}
	defer cursor.Close(ctx)

	return &result, nil
}

func (tr *TaskRepository) GetTasksByUserId(id *primitive.ObjectID) (*[]models.TaskDTO, error) {
	ctx := context.TODO()

	matchStage := bson.D{
		{"$match", bson.D{
			{
				"createdBy", id,
			},
		}},
	}

	lookupCreatedAtStage := bson.D{
		{"$lookup", bson.D{
			{
				"from",
				"users",
			},
			{
				"localField",
				"createdBy",
			},
			{
				"foreignField",
				"_id",
			},
			{
				"as",
				"createdBy",
			},
		}},
	}

	// lookupAssignedForStage := bson.D{
	// 	{"$lookup", bson.D{
	// 		{
	// 			"from",
	// 			"users",
	// 		},
	// 		{
	// 			"localField",
	// 			"assignedFor",
	// 		},
	// 		{
	// 			"foreignField",
	// 			"_id",
	// 		},
	// 		{
	// 			"as",
	// 			"assignedFor",
	// 		},
	// 	}},
	// }

	pipeline := mongo.Pipeline{matchStage, lookupCreatedAtStage}

	cursor, err := tr.db.Collection(tr.collectionName).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.TaskDTO
	for cursor.Next(ctx) {
		var taskWithUser struct {
			models.TaskDTO `bson:"inline"`
			CreatedBy      []models.UserDTO `bson:"createdBy"`
			// AssignedFor
		}

		err := cursor.Decode(&taskWithUser)
		if err != nil {
			return nil, err
		}

		if len(taskWithUser.CreatedBy) > 0 {
			taskWithUser.TaskDTO.CreatedBy = taskWithUser.CreatedBy[0]
		}

		results = append(results, taskWithUser.TaskDTO)
	}
	defer cursor.Close(ctx)

	return &results, nil
}

func (tr *TaskRepository) CreateTask(newTask *models.Task) (*mongo.InsertOneResult, error) {
	ctx := context.TODO()

	result, err := tr.db.Collection(tr.collectionName).InsertOne(ctx, newTask)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (tr *TaskRepository) UpdateTask(id *primitive.ObjectID, task *models.Task) (*mongo.UpdateResult, error) {
	isExists, err := tr.GetTaskById(id)
	if isExists == nil {
		return nil, err
	}

	ctx := context.TODO()
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"title", task.Title},
			{"description", task.Description},
			{"status", task.Status},
			{"completedAt", task.CompletedAt},
			{"updatedAt", primitive.NewDateTimeFromTime(time.Now())},
		},
		},
	}

	result, err := tr.db.Collection(tr.collectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (tr *TaskRepository) DeleteTask(id *primitive.ObjectID) (*mongo.DeleteResult, error) {
	ctx := context.TODO()
	filter := bson.D{{"_id", id}}

	result, err := tr.db.Collection(tr.collectionName).DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}
