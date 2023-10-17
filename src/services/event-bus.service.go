package services

import (
	"github.com/rpsoftech/bullion-server/src/events"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type eventBusService struct {
	eventsRepo *repos.EventRepoStruct
}

var EventBus *eventBusService

func init() {
	EventBus = &eventBusService{
		eventsRepo: repos.EventRepo,
	}
}

func (service *eventBusService) Publish(event *events.BaseEvent) {
	go service.saveToDb(event)
}
func (service *eventBusService) saveToDb(event *events.BaseEvent) {
	service.eventsRepo.Save(event)
}
