package repositories

import (
	"context"

	"github.com/kylerequez/make-you-work-app/src/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewGroupRepository(db *mongo.Database, collectionName string) *GroupRepository {
	return &GroupRepository{
		db:             db,
		collectionName: collectionName,
	}
}

func (gr *GroupRepository) GetAllGroups() (*[]models.GroupDTO, error) {
	ctx := context.TODO()

	lookupTasksStage := bson.D{{"$lookup", bson.D{
		{"from", "tasks"},
		{"localField", "tasks"},
		{"foreignField", "_id"},
		{"as", "tasks"},
	}}}

	lookupTasksCreatedByStage := bson.D{{"$lookup", bson.D{
		{"from", "users"},
		{"localField", "tasks.createdBy"},
		{"foreignField", "_id"},
		{"as", "tasksCreatedBy"},
	}}}

	lookupMembersStage := bson.D{
		{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "members"},
			{"foreignField", "_id"},
			{"as", "members"},
		}}}

	lookupCreatedByStage := bson.D{
		{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "createdBy"},
			{"foreignField", "_id"},
			{"as", "createdBy"},
		}}}

	pipeline := mongo.Pipeline{
		lookupTasksStage,
		lookupTasksCreatedByStage,
		lookupMembersStage,
		lookupCreatedByStage,
	}

	cursor, err := gr.db.Collection(gr.collectionName).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.GroupDTO
	for cursor.Next(ctx) {
		var group struct {
			models.GroupDTO `bson:"inline"`
			CreatedBy       []models.UserDTO `bson:"createdBy"`
			Tasks           []struct {
				Task      models.TaskDTO     `bson:"inline"`
				CreatedBy primitive.ObjectID `bson:"createdBy"`
			} `bson:"tasks"`
			TasksCreatedBy []models.UserDTO `bson:"tasksCreatedBy"`
			Members        []models.UserDTO `bson:"members"`
		}

		err := cursor.Decode(&group)
		if err != nil {
			return nil, err
		}

		group.GroupDTO.CreatedBy = group.CreatedBy[0]
		group.GroupDTO.Members = group.Members

		var tasks []models.TaskDTO
		for task := range group.Tasks {
			group.Tasks[task].Task.CreatedBy = group.TasksCreatedBy[task]
			tasks = append(tasks, group.Tasks[task].Task)
		}
		group.GroupDTO.Tasks = tasks

		results = append(results, group.GroupDTO)
	}

	return &results, nil
}

func (gr *GroupRepository) GetGroupById(id *primitive.ObjectID) (*models.GroupDTO, error) {
	ctx := context.TODO()

	matchStage := bson.D{
		{"$match", bson.D{
			{
				"_id", id,
			},
		}},
	}

	lookupTasksStage := bson.D{{"$lookup", bson.D{
		{"from", "tasks"},
		{"localField", "tasks"},
		{"foreignField", "_id"},
		{"as", "tasks"},
	}}}

	lookupTasksCreatedByStage := bson.D{{"$lookup", bson.D{
		{"from", "users"},
		{"localField", "tasks.createdBy"},
		{"foreignField", "_id"},
		{"as", "tasksCreatedBy"},
	}}}

	lookupMembersStage := bson.D{
		{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "members"},
			{"foreignField", "_id"},
			{"as", "members"},
		}}}

	lookupCreatedByStage := bson.D{
		{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "createdBy"},
			{"foreignField", "_id"},
			{"as", "createdBy"},
		}}}

	pipeline := mongo.Pipeline{
		matchStage,
		lookupTasksStage,
		lookupTasksCreatedByStage,
		lookupMembersStage,
		lookupCreatedByStage,
	}

	cursor, err := gr.db.Collection(gr.collectionName).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result models.GroupDTO
	for cursor.Next(ctx) {
		var group struct {
			models.GroupDTO `bson:"inline"`
			CreatedBy       []models.UserDTO `bson:"createdBy"`
			Tasks           []struct {
				Task      models.TaskDTO     `bson:"inline"`
				CreatedBy primitive.ObjectID `bson:"createdBy"`
			} `bson:"tasks"`
			TasksCreatedBy []models.UserDTO `bson:"tasksCreatedBy"`
			Members        []models.UserDTO `bson:"members"`
		}

		err := cursor.Decode(&group)
		if err != nil {
			return nil, err
		}

		group.GroupDTO.CreatedBy = group.CreatedBy[0]
		group.GroupDTO.Members = group.Members

		var tasks []models.TaskDTO
		for task := range group.Tasks {
			group.Tasks[task].Task.CreatedBy = group.TasksCreatedBy[task]
			tasks = append(tasks, group.Tasks[task].Task)
		}
		group.GroupDTO.Tasks = tasks

		result = group.GroupDTO
	}

	return &result, nil
}
