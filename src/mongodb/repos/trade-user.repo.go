package repos

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb"
	"github.com/rpsoftech/bullion-server/src/redis"
	"github.com/rpsoftech/bullion-server/src/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	TradeUserRepoStruct struct {
		collection *mongo.Collection
		redis      *redis.RedisClientStruct
	}
)

const tradeUserCollectionName = "TradeUser"

var TradeUserRepo *TradeUserRepoStruct

func init() {
	if env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
		return
	}
	coll := mongodb.MongoDatabase.Collection(tradeUserCollectionName)
	TradeUserRepo = &TradeUserRepoStruct{
		collection: coll,
		redis:      redis.InitRedisAndRedisClient(),
	}
	addUniqueIndexesToCollection([]string{"id"}, TradeUserRepo.collection)
	addIndexesToCollection([]string{"bullionId", "isActive"}, TradeUserRepo.collection)
	addComboUniqueIndexesToCollection([]string{"email", "bullionId"}, TradeUserRepo.collection)
	addComboUniqueIndexesToCollection([]string{"number", "bullionId"}, TradeUserRepo.collection)
	addComboUniqueIndexesToCollection([]string{"uNumber", "bullionId"}, TradeUserRepo.collection)
	addComboUniqueIndexesToCollection([]string{"userName", "bullionId"}, TradeUserRepo.collection)
}

func (repo *TradeUserRepoStruct) Save(entity *interfaces.TradeUserEntity) (*interfaces.TradeUserEntity, error) {
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
	go repo.cacheDataToRedis(entity)
	return entity, err
}

func (repo *TradeUserRepoStruct) FindDuplicateUser(email string, number string, userName string, bullionId string) (*[]interfaces.TradeUserEntity, error) {
	condition := bson.D{
		{
			Key: "$and",
			Value: bson.A{
				bson.D{{Key: "bullionId", Value: bullionId}},
			},
		},
		{
			Key: "$or",
			Value: bson.A{
				bson.D{{Key: "email", Value: email}},
				bson.D{{Key: "number", Value: number}},
				bson.D{{Key: "userName", Value: userName}},
			},
		},
	}
	return repo.findByFilter(&mongoDbFilter{
		conditions: &condition,
	})
}

func (repo *TradeUserRepoStruct) findByFilter(filter *mongoDbFilter) (*[]interfaces.TradeUserEntity, error) {
	var result []interfaces.TradeUserEntity
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

func (repo *TradeUserRepoStruct) FindAllInActiveUser(bullionId string) (*[]interfaces.TradeUserEntity, error) {
	return repo.findByFilter(&mongoDbFilter{
		conditions: &bson.D{
			{
				Key: "$and",
				Value: bson.A{
					bson.D{{Key: "bullionId", Value: bullionId}},
					bson.D{{Key: "isActive", Value: false}},
				},
			},
		},
	})
}

func (repo *TradeUserRepoStruct) FindOne(id string) (*interfaces.TradeUserEntity, error) {
	result := new(interfaces.TradeUserEntity)
	if redisData := repo.redis.GetStringData(fmt.Sprintf("tradeUser/%s", id)); redisData != "" {
		if err := json.Unmarshal([]byte(redisData), result); err == nil {
			result.RevertTimeStamps()
			return result, err
		}
	}
	err := repo.collection.FindOne(mongodb.MongoCtx, bson.D{{
		Key: "id", Value: id,
	}}).Decode(result)

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
	go repo.cacheDataToRedis(result)
	return result, err
}

func (repo *TradeUserRepoStruct) cacheDataToRedis(entity *interfaces.TradeUserEntity) {
	entity.AddTimeStamps()
	if entityStringBytes, err := json.Marshal(entity); err == nil {
		entityString := string(entityStringBytes)
		repo.redis.SetStringDataWithExpiry(fmt.Sprintf("tradeUser/%s", entity.ID), entityString, time.Duration(24)*time.Hour)
	}
}
func (repo *TradeUserRepoStruct) findOneByCondition(bullionId string, condition *bson.E) (*interfaces.TradeUserEntity, error) {
	var result interfaces.TradeUserEntity

	err := repo.collection.FindOne(mongodb.MongoCtx, bson.D{
		{Key: "bullionId", Value: bullionId},
		*condition,
	}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// This error means your query did not match any documents.
			err = &interfaces.RequestError{
				StatusCode: http.StatusBadRequest,
				Code:       interfaces.ERROR_ENTITY_NOT_FOUND,
				Message:    fmt.Sprintf("Feeds Entity identified by %s %s not found", condition.Key, condition.Value),
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

func (repo *TradeUserRepoStruct) FindOneByEmail(bullionId string, email string) (*interfaces.TradeUserEntity, error) {
	return repo.findOneByCondition(bullionId, &bson.E{Key: "email", Value: email})
}

func (repo *TradeUserRepoStruct) FindOneByNumber(bullionId string, number string) (*interfaces.TradeUserEntity, error) {
	return repo.findOneByCondition(bullionId, &bson.E{Key: "number", Value: number})
}

func (repo *TradeUserRepoStruct) FindOneByUNumber(bullionId string, uNumber string) (*interfaces.TradeUserEntity, error) {
	return repo.findOneByCondition(bullionId, &bson.E{Key: "uNumber", Value: uNumber})
}
