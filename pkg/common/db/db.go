package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Client {
	mongoURI := os.Getenv("MONGO_URI")

	if mongoURI == "" {
		log.Fatal("You must set MONGO_URI environment variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))

	if err != nil {
		log.Fatal("Error connecting database", err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	dbName := os.Getenv("DB_NAME")
	database := client.Database(dbName)

	_, err = database.Collection("users").Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected")

	return client
}
