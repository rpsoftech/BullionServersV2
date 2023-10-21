package repos

import (
	"errors"
	"fmt"

	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FeedsRepoStruct struct {
	collection *mongo.Collection
}

const feedCollectionName = "Feed"

var FeedsRepo *FeedsRepoStruct

func init() {
	if env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
		return
	}
	coll := mongodb.MongoDatabase.Collection(feedCollectionName)
	FeedsRepo = &FeedsRepoStruct{
		collection: coll,
	}
	addUniqueIndexesToCollection([]string{"id"}, FeedsRepo.collection)
	addIndexesToCollection([]string{"bullionId", "createdAt"}, FeedsRepo.collection)
}

func (repo *FeedsRepoStruct) Save(entity *interfaces.FeedsEntity) error {
	_, err := repo.collection.InsertOne(mongodb.MongoCtx, entity)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			err = &interfaces.RequestError{
				StatusCode: 500,
				Code:       interfaces.ERROR_INTERNAL_SERVER,
				Message:    fmt.Sprintf("Internal Server Error: %s", err.Error()),
				Name:       "INTERNAL_ERROR",
			}
		} else {
			err = nil
		}
	}
	return err
}

func (repo *FeedsRepoStruct) findByFilter(filter bson.D) (*[]interfaces.FeedsEntity, error) {
	var result []interfaces.FeedsEntity
	cursor, err := repo.collection.Find(mongodb.MongoCtx, filter)
	if err == nil {
		err = cursor.All(mongodb.MongoCtx, &result)
	}
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// This error means your query did not match any documents.
			err = &interfaces.RequestError{
				StatusCode: 400,
				Code:       interfaces.ERROR_ENTITY_NOT_FOUND,
				Message:    fmt.Sprintf("Feeds Entities filtered By %v not found", filter),
				Name:       "ENTITY_NOT_FOUND",
			}
		} else {
			err = &interfaces.RequestError{
				StatusCode: 500,
				Code:       interfaces.ERROR_INTERNAL_SERVER,
				Message:    fmt.Sprintf("Internal Server Error: %s", err.Error()),
				Name:       "INTERNAL_ERROR",
			}
		}
	}
	return &result, err
}

func (repo *FeedsRepoStruct) GetAllByBullionId(bullionId string) (*[]interfaces.FeedsEntity, error) {
	return repo.findByFilter(bson.D{{Key: "bullionId", Value: bullionId}})
}
