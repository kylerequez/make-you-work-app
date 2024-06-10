package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupDTO struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	CreatedBy   UserDTO            `bson:"createdBy" json:"createdBy"`
	CreatedAt   primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedAt   primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	Tasks       []TaskDTO          `bson:"tasks" json:"tasks"`
	Members     []UserDTO          `bson:"members" json:"members"`
}
