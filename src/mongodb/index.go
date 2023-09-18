package mongodb

import (
	"context"

	"github.com/rpsoftech/bullion-server/src/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database
var MongoCtx = context.TODO()

func init() {
	// env.Env.DB_URL
	client, err := mongo.Connect(MongoCtx, options.Client().ApplyURI(env.Env.DB_URL))
	if err != nil {
		panic(err)
	}
	MongoClient = client
	MongoDatabase = client.Database(env.Env.DB_NAME)
}

func DeferFunction() {
	if err := MongoClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
