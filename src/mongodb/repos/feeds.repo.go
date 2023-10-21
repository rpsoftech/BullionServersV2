package repos

import (
	"errors"
	"fmt"

	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb"
	"github.com/rpsoftech/bullion-server/src/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	FeedsRepoStruct struct {
		collection *mongo.Collection
	}

	feedFilters struct {
		conditions *bson.D
		sort       *bson.D
		limit      int64
		skip       int64
	}
)

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

func (repo *FeedsRepoStruct) Save(entity *interfaces.FeedsEntity) (*interfaces.FeedsEntity, error) {
	if err := utility.ValidateStructAndReturnReqError(entity, &interfaces.RequestError{
		StatusCode: 400,
		Code:       interfaces.ERROR_INVALID_ENTITY,
		Message:    "",
		Name:       "ERROR_INVALID_ENTITY",
	}); err != nil {
		return entity, err
	}
	entity.Updated()
	err := repo.collection.FindOneAndUpdate(mongodb.MongoCtx, bson.D{{
		Key: "_id", Value: entity.ID,
	}}, bson.D{{Key: "$set", Value: entity}}, findOneAndUpdateOptions).Err()
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
	return entity, err
}

func (repo *FeedsRepoStruct) findByFilter(filter *feedFilters) (*[]interfaces.FeedsEntity, error) {
	var result []interfaces.FeedsEntity
	opt := options.Find()
	if filter.sort != nil {
		opt.SetSort(filter.sort)
	}
	if filter.limit > 0 {
		opt.SetLimit(filter.limit)
	}
	if filter.skip > 0 {
		opt.SetSkip(filter.skip)
	}
	cursor, err := repo.collection.Find(mongodb.MongoCtx, filter.conditions, opt)
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
	return repo.findByFilter(&feedFilters{
		conditions: &bson.D{{Key: "bullionId", Value: bullionId}},
	})
}

func (repo *FeedsRepoStruct) DeleteById(id string) error {
	_, err := repo.collection.DeleteOne(mongodb.MongoCtx, bson.D{{
		Key: "id", Value: id,
	}})

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// This error means your query did not match any documents.
			err = &interfaces.RequestError{
				StatusCode: 400,
				Code:       interfaces.ERROR_ENTITY_NOT_FOUND,
				Message:    fmt.Sprintf("Feeds Entity identified by id %s not found", id),
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
	return err
}

func (repo *FeedsRepoStruct) GetPaginatedFeedInDescendingOrder(bullionId string, page int64, limit int64) (*[]interfaces.FeedsEntity, error) {
	println(limit)
	return repo.findByFilter(&feedFilters{
		conditions: &bson.D{{Key: "bullionId", Value: bullionId}},
		sort:       &bson.D{{Key: "createdAt", Value: -1}},
		limit:      limit,
		skip:       page * limit,
	})
}

func (repo *FeedsRepoStruct) FindOne(id string) (*interfaces.FeedsEntity, error) {
	var result interfaces.FeedsEntity

	err := repo.collection.FindOne(mongodb.MongoCtx, bson.D{{
		Key: "id", Value: id,
	}}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// This error means your query did not match any documents.
			err = &interfaces.RequestError{
				StatusCode: 400,
				Code:       interfaces.ERROR_ENTITY_NOT_FOUND,
				Message:    fmt.Sprintf("Feeds Entity identified by id %s not found", id),
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
