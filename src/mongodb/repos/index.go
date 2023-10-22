package repos

import (
	"github.com/rpsoftech/bullion-server/src/mongodb"
	"github.com/rpsoftech/bullion-server/src/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDbFilter struct {
	conditions *bson.D
	sort       *bson.D
	limit      int64
	skip       int64
}

var findOneAndUpdateOptions = &options.FindOneAndUpdateOptions{
	Upsert: utility.BoolPointer(true),
}

func addComboUniqueIndexesToCollection(UniqueIndexes []string, collection *mongo.Collection) {
	i := bson.D{}
	for _, element := range UniqueIndexes {
		i = append(i, bson.E{Key: element, Value: 1})
	}
	collection.Indexes().CreateOne(mongodb.MongoCtx, mongo.IndexModel{
		Keys:    i,
		Options: options.Index().SetUnique(true),
	})
}

func addIndexesToCollection(indexesString []string, collection *mongo.Collection) {
	indexArray := make([]mongo.IndexModel, len(indexesString))
	for i, element := range indexesString {
		indexArray[i] = mongo.IndexModel{
			Keys:    bson.D{{Key: element, Value: 1}},
			Options: options.Index().SetUnique(false),
		}
	}
	collection.Indexes().CreateMany(mongodb.MongoCtx, indexArray)
}
func addUniqueIndexesToCollection(UniqueIndexes []string, collection *mongo.Collection) {
	indexes := make([]mongo.IndexModel, len(UniqueIndexes))
	for i, element := range UniqueIndexes {
		indexes[i] = mongo.IndexModel{
			Keys:    bson.D{{Key: element, Value: 1}},
			Options: options.Index().SetUnique(true),
		}
	}
	collection.Indexes().CreateMany(mongodb.MongoCtx, indexes)
}
