package db

import (
	"context"
	"fmt"
	"log"

	"github.com/kylerequez/make-you-work-app/src/api/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

type DBCredentials struct {
	host string
	port string
	uri  string
}

func ConnectDB() error {
	log.Println(":::-::: Connecting to DB...")
	credentials, err := GetDBCredentials()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf(credentials.uri, credentials.host, credentials.port)

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return err
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

func GetDBCredentials() (*DBCredentials, error) {
	host, err := utils.GetEnv("DB_HOST")
	if err != nil {
		return nil, err
	}
	port, err := utils.GetEnv("DB_PORT")
	if err != nil {
		return nil, err
	}
	uri, err := utils.GetEnv("DB_URI_TEMPLATE")
	if err != nil {
		return nil, err
	}

	return &DBCredentials{
		host: *host,
		port: *port,
		uri:  *uri,
	}, nil
}
