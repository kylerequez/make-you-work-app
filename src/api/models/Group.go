package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Group struct {
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description" json:"description"`
	CreatedBy   primitive.ObjectID   `bson:"createdBy" json:"createdBy"`
	CreatedAt   primitive.DateTime   `bson:"createdBy" json:"createdAt"`
	UpdatedAt   primitive.DateTime   `bson:"updatedAt" json:"updatedAt"`
	Tasks       []primitive.ObjectID `bson:"tasks" json:"tasks"`
	Members     []primitive.ObjectID `bson:"members" json:"members"`
}
