package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect() *mongo.Database {
	connectionUrl := os.Getenv("CONNECTION_URL")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionUrl))
	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	databaseName := os.Getenv("DATABASE_NAME")
	database := client.Database(databaseName)

	return database
}
