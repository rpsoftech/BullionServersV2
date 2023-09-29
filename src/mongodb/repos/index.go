package repos

import (
	"github.com/rpsoftech/bullion-server/src/mongodb"
	"github.com/rpsoftech/bullion-server/src/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var findOneAndUpdateOptions = &options.FindOneAndUpdateOptions{
	Upsert: utility.BoolPointer(true),
}

func addUniqueIndexesToCollection(UniqueIndexes []string, collection *mongo.Collection) {
	for _, element := range UniqueIndexes {
		collection.Indexes().CreateOne(mongodb.MongoCtx, mongo.IndexModel{
			Keys:    bson.D{{Key: element, Value: 1}},
			Options: options.Index().SetUnique(true),
		})
	}
}
