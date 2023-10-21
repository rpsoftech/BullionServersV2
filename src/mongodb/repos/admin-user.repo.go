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

type AdminUserRepoStruct struct {
	collection *mongo.Collection
}

const adminUserCollectionName = "AdminUser"

var AdminUserRepo *AdminUserRepoStruct

func init() {
	if env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
		return
	}
	coll := mongodb.MongoDatabase.Collection(adminUserCollectionName)
	AdminUserRepo = &AdminUserRepoStruct{
		collection: coll,
	}
	addUniqueIndexesToCollection([]string{"id"}, AdminUserRepo.collection)
	addComboUniqueIndexesToCollection([]string{"userName", "bullionId"}, AdminUserRepo.collection)
}

func (repo *AdminUserRepoStruct) Save(entity *interfaces.AdminUserEntity) (*interfaces.AdminUserEntity, error) {
	var result interfaces.AdminUserEntity
	err := repo.collection.FindOneAndUpdate(mongodb.MongoCtx, bson.D{{
		Key: "_id", Value: entity.ID,
	}}, bson.D{{Key: "$set", Value: entity}}, findOneAndUpdateOptions).Decode(&result)
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
	return &result, err
}

func (repo *AdminUserRepoStruct) FindOne(id string) (*interfaces.AdminUserEntity, error) {
	var result interfaces.AdminUserEntity
	err := repo.collection.FindOne(mongodb.MongoCtx, bson.D{{
		Key: "id", Value: id,
	}}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// This error means your query did not match any documents.
			err = &interfaces.RequestError{
				StatusCode: 400,
				Code:       interfaces.ERROR_ENTITY_NOT_FOUND,
				Message:    fmt.Sprintf("GeneralUser Entity identified by id %s not found", id),
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

func (repo *AdminUserRepoStruct) FindOneUserNameAndBullionId(uname string, bullionId string) (*interfaces.AdminUserEntity, error) {
	var result interfaces.AdminUserEntity
	err := repo.collection.FindOne(mongodb.MongoCtx, bson.D{{
		Key: "userName", Value: uname,
	}, {
		Key: "bullionId", Value: bullionId,
	}}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// This error means your query did not match any documents.
			err = &interfaces.RequestError{
				StatusCode: 400,
				Code:       interfaces.ERROR_ENTITY_NOT_FOUND,
				Message:    fmt.Sprintf("GeneralUser Entity identified by uname %s and bullionId %s not found", uname, bullionId),
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
