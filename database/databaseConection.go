package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = dbDatabase()

const connectionString = "mongodb://localhost:27017"

func dbDatabase() *mongo.Client {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Could not make connection with database, ", err)
	}
	fmt.Println("Connected to MongoDB!!")
	return client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("bill").Collection(collectionName)
	return collection
}
