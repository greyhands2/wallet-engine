package config

import (
	"context"
	"fmt"
	"log"
	"github.com/greyhands2/wallet-engine/vipr"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)
var connectionString string = vipr.ViperEnvVariable("MONGO_URL")
var dbName string = vipr.ViperEnvVariable("MONGO_DB")
func connect() *mongo.Client {
	//client option

	
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
	
	var collection *mongo.Collection = client.Database(dbName).Collection(collectionName)

	return collection
}
