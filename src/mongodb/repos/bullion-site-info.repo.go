package repos

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BullionSiteInfoRepoStruct struct {
	collection *mongo.Collection
}

const bullionSiteInfoCollectionName = "BullionSiteInfo"

var BullionSiteInfoRepo *BullionSiteInfoRepoStruct

func init() {
	if env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
		return
	}
	coll := mongodb.MongoDatabase.Collection(bullionSiteInfoCollectionName)
	BullionSiteInfoRepo = &BullionSiteInfoRepoStruct{
		collection: coll,
	}
	addUniqueIndexesToCollection([]string{"id", "domain", "shortName"}, BullionSiteInfoRepo.collection)
}

func (repo *BullionSiteInfoRepoStruct) Save(entity *interfaces.BullionSiteInfoEntity) (*interfaces.BullionSiteInfoEntity, error) {
	var result interfaces.BullionSiteInfoEntity
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

func (repo *BullionSiteInfoRepoStruct) FindOne(id string) (*interfaces.BullionSiteInfoEntity, error) {
	var result interfaces.BullionSiteInfoEntity
	err := repo.collection.FindOne(mongodb.MongoCtx, bson.D{{
		Key: "id", Value: id,
	}}).Decode(&result)

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
	return &result, err
}

func (repo *BullionSiteInfoRepoStruct) FindOneByDomain(domain string) (*interfaces.BullionSiteInfoEntity, error) {
	var result interfaces.BullionSiteInfoEntity

	err := repo.collection.FindOne(mongodb.MongoCtx, bson.D{{
		Key: "domain", Value: domain,
	}}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// This error means your query did not match any documents.
			err = &interfaces.RequestError{
				StatusCode: http.StatusBadRequest,
				Code:       interfaces.ERROR_ENTITY_NOT_FOUND,
				Message:    fmt.Sprintf("Bullion Entity identified by domain %s not found", domain),
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

func (repo *BullionSiteInfoRepoStruct) FindByShortName(name string) (*interfaces.BullionSiteInfoEntity, error) {
	var result interfaces.BullionSiteInfoEntity
	err := repo.collection.FindOne(mongodb.MongoCtx, bson.D{{
		Key: "shortName", Value: name,
	}}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// This error means your query did not match any documents.
			err = &interfaces.RequestError{
				StatusCode: http.StatusBadRequest,
				Code:       interfaces.ERROR_ENTITY_NOT_FOUND,
				Message:    fmt.Sprintf("Bullion Entity identified by shortname %s not found", name),
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
