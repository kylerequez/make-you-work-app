package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname   string             `json:"firstname" bson:"firstname"`
	Middlename  string             `json:"middlename" bson:"middlename"`
	Lastname    string             `json:"lastname" bson:"lastname"`
	Authorities []string           `json:"authorities" bson:"authorities"`
	Status      string             `json:"status" bson:"status"`
	Email       string             `json:"email" bson:"email"`
	Username    string             `json:"username" bson:"username"`
	Password    []byte             `json:"password" bson:"password"`
	CreatedAt   primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt   primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}
