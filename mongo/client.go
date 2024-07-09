package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
)

func CreateMongoClient() *mongo.Client {
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().SetAppName("assets-scheduler").ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	return client
}
