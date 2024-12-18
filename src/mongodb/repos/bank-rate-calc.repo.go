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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BankRateCalcRepoStruct struct {
	collection *mongo.Collection
	redis      *redis.RedisClientStruct
	// historyCollection *mongo.Collection
}

const bankRateCalcRepoCollectionName = "BankRateCalc"
const bankRateRedisCollection = "bankRate"

// const bankRateCalcHistoryRepoCollectionName = "BankRateCalcHistory"

var BankRateCalcRepo *BankRateCalcRepoStruct

func init() {
	if env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
		return
	}
	coll := mongodb.MongoDatabase.Collection(bankRateCalcRepoCollectionName)
	BankRateCalcRepo = &BankRateCalcRepoStruct{
		collection: coll,
		redis:      redis.InitRedisAndRedisClient(),
	}
	addUniqueIndexesToCollection([]string{"id", "bullionId"}, BankRateCalcRepo.collection)
}

func (repo *BankRateCalcRepoStruct) cacheDataToRedis(entity *interfaces.BankRateCalcEntity) {
	if entityStringBytes, err := json.Marshal(entity); err == nil {
		entityString := string(entityStringBytes)
		repo.redis.SetStringDataWithExpiry(fmt.Sprintf("%s/%s", bankRateRedisCollection, entity.BullionId), entityString, time.Duration(24)*time.Hour)
	}
}

func (repo *BankRateCalcRepoStruct) Save(entity *interfaces.BankRateCalcEntity) (*interfaces.BankRateCalcEntity, error) {
	var result interfaces.BankRateCalcEntity
	err := repo.collection.FindOneAndUpdate(mongodb.MongoCtx, bson.D{{
		Key: "_id", Value: entity.ID,
	}}, bson.D{{Key: "$set", Value: entity}}, findOneAndUpdateOptions).Decode(&result)
	entity.Updated()
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
	return &result, err
}

func (repo *BankRateCalcRepoStruct) FindOneByBullionId(id string) (*interfaces.BankRateCalcEntity, error) {
	result := new(interfaces.BankRateCalcEntity)
	if redisData := repo.redis.GetStringData(fmt.Sprintf("%s/%s", bankRateRedisCollection, id)); redisData != "" {
		if err := json.Unmarshal([]byte(redisData), result); err == nil {
			return result, err
		}
	}
	err := repo.collection.FindOne(mongodb.MongoCtx, bson.D{{
		Key: "bullionId", Value: id,
	}}).Decode(result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// This error means your query did not match any documents.
			err = &interfaces.RequestError{
				StatusCode: http.StatusBadRequest,
				Code:       interfaces.ERROR_ENTITY_NOT_FOUND,
				Message:    fmt.Sprintf("Bullion Entity identified by id %s not found", id),
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