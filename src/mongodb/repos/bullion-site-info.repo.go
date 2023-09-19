package repos

import (
	"github.com/rpsoftech/bullion-server/src/interfaces/bullion"
	"github.com/rpsoftech/bullion-server/src/mongodb"
	"github.com/rpsoftech/bullion-server/src/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const bullionSiteInfoCollectionName = "BullionSiteInfo"

var (
	BullionSiteInfoRepo     *bullionSiteInfoRepo
	findOneAndUpdateOptions = &options.FindOneAndUpdateOptions{
		Upsert: utility.BoolPointer(true),
	}
)

func init() {
	coll := mongodb.MongoDatabase.Collection(bullionSiteInfoCollectionName)
	BullionSiteInfoRepo = &bullionSiteInfoRepo{
		collection: coll,
	}
	BullionSiteInfoRepo.collection.Indexes().CreateOne(mongodb.MongoCtx, mongo.IndexModel{
		// Options.Background: utility.BoolPointer(true),
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
}

type bullionSiteInfoRepo struct {
	collection *mongo.Collection
}

func (repo *bullionSiteInfoRepo) Save(entity *bullion.BullionSiteInfo) (result *bullion.BullionSiteInfo, err error) {
	err = repo.collection.FindOneAndUpdate(mongodb.MongoCtx, bson.D{{
		Key: "_id", Value: entity.ID,
	}}, bson.D{{Key: "$set", Value: entity}}, findOneAndUpdateOptions).Decode(result)
	return
}

func (repo *bullionSiteInfoRepo) FindOne(id string) (result *bullion.BullionSiteInfo) {
	repo.collection.FindOne(mongodb.MongoCtx, bson.D{{
		Key: "_id", Value: id,
	}}).Decode(result)
	return
}
