package repos

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb"
	"github.com/rpsoftech/bullion-server/src/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	ProductGroupMapRepoStruct struct {
		collection *mongo.Collection
	}
)

const productGroupMapCollectionName = "ProductGroupMap"

var ProductGroupMapRepo *ProductGroupMapRepoStruct

func init() {
	if env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
		return
	}
	coll := mongodb.MongoDatabase.Collection(productGroupMapCollectionName)
	ProductGroupMapRepo = &ProductGroupMapRepoStruct{
		collection: coll,
	}

	addUniqueIndexesToCollection([]string{"id"}, ProductGroupMapRepo.collection)
	addComboUniqueIndexesToCollection([]string{"groupId", "productId"}, ProductGroupMapRepo.collection)
	addComboIndexesToCollection([]string{"bullionId", "groupId"}, ProductGroupMapRepo.collection)
	addIndexesToCollection([]string{"bullionId", "createdAt", "groupId"}, ProductGroupMapRepo.collection)
}

func (repo *ProductGroupMapRepoStruct) Save(entity *interfaces.TradeUserGroupMapEntity) (*interfaces.TradeUserGroupMapEntity, error) {
	if err := utility.ValidateStructAndReturnReqError(entity, &interfaces.RequestError{
		StatusCode: http.StatusBadRequest,
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

func (repo *ProductGroupMapRepoStruct) BulkUpdate(entities *[]interfaces.TradeUserGroupMapEntity) (*[]interfaces.TradeUserGroupMapEntity, error) {
	models := make([]mongo.WriteModel, len(*entities))
	for i, entity := range *entities {
		if err := utility.ValidateStructAndReturnReqError(entity, &interfaces.RequestError{
			StatusCode: http.StatusBadRequest,
			Code:       interfaces.ERROR_INVALID_ENTITY,
			Message:    "",
			Name:       "ERROR_INVALID_ENTITY",
		}); err != nil {
			return nil, err
		}
		entity.Updated()
		models[i] = mongo.NewUpdateOneModel().SetFilter(bson.D{{Key: "_id", Value: entity.ID}}).SetUpdate(
			bson.D{{Key: "$set", Value: entity}}).SetUpsert(true)
	}
	_, err := repo.collection.BulkWrite(mongodb.MongoCtx, models)
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
	return entities, err
}

func (repo *ProductGroupMapRepoStruct) findByFilter(filter *mongoDbFilter) (*[]interfaces.TradeUserGroupMapEntity, error) {
	var result []interfaces.TradeUserGroupMapEntity
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
				StatusCode: http.StatusBadRequest,
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

func (repo *ProductGroupMapRepoStruct) GetAllByBullionId(bullionId string) (*[]interfaces.TradeUserGroupMapEntity, error) {
	return repo.findByFilter(&mongoDbFilter{
		conditions: &bson.D{{Key: "bullionId", Value: bullionId}},
	})
}
func (repo *ProductGroupMapRepoStruct) GetAllByGroupId(groupId string, bullionId string) (*[]interfaces.TradeUserGroupMapEntity, error) {
	return repo.findByFilter(&mongoDbFilter{
		conditions: &bson.D{
			{Key: "groupId", Value: groupId},
			{Key: "bullionId", Value: bullionId},
		},
	})
}

func (repo *ProductGroupMapRepoStruct) FindOne(id string) (*interfaces.TradeUserGroupMapEntity, error) {
	var result interfaces.TradeUserGroupMapEntity
	err := repo.collection.FindOne(mongodb.MongoCtx, bson.D{{
		Key: "id", Value: id,
	}}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// This error means your query did not match any documents.
			err = &interfaces.RequestError{
				StatusCode: http.StatusBadRequest,
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
