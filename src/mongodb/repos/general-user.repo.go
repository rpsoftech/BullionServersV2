package repos

import (
	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type generalUserRepo struct {
	collection *mongo.Collection
}

const generalUserCollectionName = "GeneralUser"

var GeneralUserRepo *generalUserRepo

func init() {
	if env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
		return
	}
	coll := mongodb.MongoDatabase.Collection(generalUserCollectionName)
	GeneralUserRepo = &generalUserRepo{
		collection: coll,
	}
	addUniqueIndexesToCollection([]string{"id"}, GeneralUserRepo.collection)
}

func (repo *generalUserRepo) Save(entity *interfaces.GeneralUserEntity) (result interfaces.GeneralUserEntity, err error) {
	err = repo.collection.FindOneAndUpdate(mongodb.MongoCtx, bson.D{{
		Key: "_id", Value: entity.ID,
	}}, bson.D{{Key: "$set", Value: entity}}, findOneAndUpdateOptions).Decode(&result)
	return
}

func (repo *generalUserRepo) FindOne(id string) (result interfaces.GeneralUserEntity) {
	repo.collection.FindOne(mongodb.MongoCtx, bson.D{{
		Key: "id", Value: id,
	}}).Decode(&result)
	return
}
