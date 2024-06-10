package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

func ConnectDB() error {
	log.Println(":::-::: Connecting to DB...")
	var dbHost, dbPort, dbString = os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_URI_TEMPLATE")

	var dbURI string = fmt.Sprintf(dbString, dbHost, dbPort)

	clientOpts := options.Client().ApplyURI(dbURI)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return errors.New(err.Error())
	}

	dbClient = client
	log.Println(":::-::: Successfully Connected to DB")
	return nil
}

func CloseDB() error {
	log.Println(":::-::: Closed DB")
	return dbClient.Disconnect(context.TODO())
}

func GetDB(database string) *mongo.Database {
	return dbClient.Database(database)
}
