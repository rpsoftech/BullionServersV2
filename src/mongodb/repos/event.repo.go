package repos

import (
	"fmt"

	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/events"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventRepoStruct struct {
	collection *mongo.Collection
}

const eventsCollectionName = "Events"

var EventRepo *EventRepoStruct

func init() {
	if env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
		return
	}
	coll := mongodb.MongoDatabase.Collection(eventsCollectionName)
	EventRepo = &EventRepoStruct{
		collection: coll,
	}
	addUniqueIndexesToCollection([]string{"id"}, EventRepo.collection)
}

func (repo *EventRepoStruct) Save(entity *events.BaseEvent) error {
	_, err := repo.collection.InsertOne(mongodb.MongoCtx, entity)
	if err != nil {
		if err != mongo.ErrNoDocuments {
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
