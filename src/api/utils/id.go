package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HexToObjectId(id string) (*primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return &oid, nil
}
