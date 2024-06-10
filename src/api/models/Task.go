package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	Title       string              `bson:"title" json:"title"`
	Description string              `bson:"description" json:"description"`
	Status      string              `bson:"status" json:"status"`
	CompletedAt *primitive.DateTime `bson:"completedAt" json:"completedAt"`
	CreatedBy   *primitive.ObjectID `json:"createdBy" bson:"createdBy"`
	CreatedAt   primitive.DateTime  `json:"createdAt" bson:"createdAt"`
	UpdatedAt   primitive.DateTime  `json:"updatedAt" bson:"updatedAt"`
}
