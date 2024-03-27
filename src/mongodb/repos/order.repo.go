package repos

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb"
	"github.com/rpsoftech/bullion-server/src/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	OrderRepoStruct struct {
		collection *mongo.Collection
	}
)

const orderCollectionName = "Order"

var OrderRepo *OrderRepoStruct

func init() {
	if env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
		return
	}
	coll := mongodb.MongoDatabase.Collection(orderCollectionName)
	OrderRepo = &OrderRepoStruct{
		collection: coll,
	}
	addUniqueIndexesToCollection([]string{"id"}, OrderRepo.collection)
	addIndexesToCollection([]string{"userId", "productGroupMapId", "groupId", "productId", "orderStatus", "createdAt"}, OrderRepo.collection)
}

func (repo *OrderRepoStruct) Save(entity *interfaces.OrderEntity) (*interfaces.OrderEntity, error) {
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

func (repo *OrderRepoStruct) BulkUpdate(entities *[]interfaces.OrderEntity) (*[]interfaces.OrderEntity, error) {
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

func (repo *OrderRepoStruct) findByFilter(filter *mongoDbFilter) (*[]interfaces.OrderEntity, error) {
	var result []interfaces.OrderEntity
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

func (repo *OrderRepoStruct) FindOne(id string) (*interfaces.OrderEntity, error) {
	var result interfaces.OrderEntity

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

func (repo *OrderRepoStruct) DeleteOrderHistoryById(id string) error {
	_, err := repo.collection.DeleteOne(mongodb.MongoCtx, bson.D{{
		Key: "_id", Value: id,
	}})
	return err
}

func (repo *OrderRepoStruct) GetOrdersByBullionIdWithDateRangeAndOrderStatus(bullionId string, startDate time.Time, endDate time.Time, orderStatusArray *[]interfaces.OrderStatus) (*[]interfaces.OrderEntity, error) {
	return repo.findByFilter(&mongoDbFilter{
		sort: &bson.D{{Key: "createdAt", Value: -1}},
		conditions: &bson.D{
			{Key: "bullionId", Value: bullionId},
			{Key: "createdAt", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}},
			{Key: "orderStatus", Value: bson.D{{Key: "$in", Value: *orderStatusArray}}},
		},
	})
}

func (repo *OrderRepoStruct) GetUsersOrderPaginated(userId string, page int64, limit int64) (*[]interfaces.OrderEntity, error) {
	return repo.findByFilter(&mongoDbFilter{
		sort: &bson.D{{Key: "createdAt", Value: -1}},
		conditions: &bson.D{
			{Key: "userId", Value: userId},
		},
		limit: limit,
		skip:  page * limit,
	})
}

func (repo *OrderRepoStruct) GetUsersOrderPaginatedWithOrderStatusArray(userId string, orderStatusArray *[]interfaces.OrderStatus, page int64, limit int64) (*[]interfaces.OrderEntity, error) {
	return repo.findByFilter(&mongoDbFilter{
		sort: &bson.D{{Key: "createdAt", Value: -1}},
		conditions: &bson.D{
			{Key: "userId", Value: userId},
			{Key: "orderStatus", Value: bson.D{{Key: "$in", Value: *orderStatusArray}}},
		},
		limit: limit,
		skip:  page * limit,
	})
}
