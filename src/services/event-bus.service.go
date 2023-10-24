package services

import (
	"github.com/rpsoftech/bullion-server/src/events"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type eventBusService struct {
	eventsRepo *repos.EventRepoStruct
}

var eventBus *eventBusService

func getEventBusService() *eventBusService {
	if eventBus == nil {
		eventBus = &eventBusService{
			eventsRepo: repos.EventRepo,
		}
		println("EventBus Service Initialized")
	}
	return eventBus
}
func (service *eventBusService) Publish(event *events.BaseEvent) {
	go service.saveToDb(event)
}
func (service *eventBusService) PublishAll(event *[]interface{}) {
	go service.saveAllToDb(event)
}
func (service *eventBusService) saveAllToDb(events *[]interface{}) {
	service.eventsRepo.SaveAll(events)
}
func (service *eventBusService) saveToDb(event *events.BaseEvent) {
	service.eventsRepo.Save(event)
}
