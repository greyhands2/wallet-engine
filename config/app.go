package config

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func connect() *mongo.Client {
	//client option

	connectionString := "mongodb+srv://starq:Test1234@cluster0.phlp6.mongodb.net/test"
	fmt.Printf("the url is %s\n", connectionString)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	return client
}

var Client *mongo.Client = connect()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {

	var collection *mongo.Collection = client.Database("wallet-engine").Collection(collectionName)

	return collection
}
