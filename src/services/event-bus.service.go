package services

import (
	"github.com/rpsoftech/bullion-server/src/events"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
	"github.com/rpsoftech/bullion-server/src/redis"
)

type eventBusService struct {
	eventsRepo *repos.EventRepoStruct
	redis      *redis.RedisClientStruct
}

var eventBus *eventBusService

func getEventBusService() *eventBusService {
	if eventBus == nil {
		eventBus = &eventBusService{
			eventsRepo: repos.EventRepo,
			redis:      redis.InitRedisAndRedisClient(),
		}
		println("EventBus Service Initialized")
	}
	return eventBus
}
func (service *eventBusService) Publish(event *events.BaseEvent) {
	go service.saveToDb(event)
}
func (service *eventBusService) PublishAll(event *[]events.BaseEvent) {
	go service.saveAllToDb(event)
}
func (service *eventBusService) saveAllToDb(events *[]events.BaseEvent) {
	for _, event := range *events {
		service.saveToDb(&event)
	}
}
func (service *eventBusService) saveToDb(event *events.BaseEvent) {
	service.redis.PublishEvent(event)
	service.eventsRepo.Save(event)
}
