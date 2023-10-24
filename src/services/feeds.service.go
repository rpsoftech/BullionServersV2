package services

import (
	"github.com/rpsoftech/bullion-server/src/events"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type feedsService struct {
	feedsRepo  *repos.FeedsRepoStruct
	eventBus   *eventBusService
	fcmService *firebaseCloudMessagingService
}

var FeedsService *feedsService

func init() {
	FeedsService = &feedsService{
		feedsRepo:  repos.FeedsRepo,
		eventBus:   getEventBusService(),
		fcmService: getFirebaseCloudMessagingService(),
	}
	println("Feed Service Initialized")
}

func (service *feedsService) SendNotification(bullionId string, entity *interfaces.FeedsBase, adminId string) error {

	service.fcmService.SendTextNotificationToAll(bullionId, entity.Title, entity.Body, entity.IsHtml)
	event := events.CreateNotificationSendEvent(entity, adminId)
	service.eventBus.Publish(event)
	return nil
}

func (service *feedsService) UpdateFeeds(baseEntity *interfaces.FeedsBase, feedId string, adminId string) (*interfaces.FeedsEntity, error) {
	entity, err := service.feedsRepo.FindOne(feedId)
	if err != nil {
		return nil, err
	}
	if entity.BullionId != baseEntity.BullionId {
		return nil, &interfaces.RequestError{
			StatusCode: 403,
			Code:       interfaces.ERROR_MISMATCH_BULLION_ID,
			Message:    "You do not have access to this Feed",
			Name:       "ERROR_MISMATCH_BULLION_ID",
		}
	}
	entity.FeedsBase = baseEntity
	entity.Updated()
	return service.AddAndUpdateNewFeeds(entity, adminId)
}

func (service *feedsService) AddAndUpdateNewFeeds(entity *interfaces.FeedsEntity, adminID string) (*interfaces.FeedsEntity, error) {
	event := events.CreateUpdateFeedEvent(entity, adminID)
	go service.eventBus.Publish(event)
	return service.feedsRepo.Save(entity)
}

func (service *feedsService) FetchAllFeedsByBullionId(bullionId string) (*[]interfaces.FeedsEntity, error) {
	return service.feedsRepo.GetAllByBullionId(bullionId)
}

func (service *feedsService) FetchPaginatedFeedsByBullionId(bullionId string, page int64, limit int64) (*[]interfaces.FeedsEntity, error) {
	return service.feedsRepo.GetPaginatedFeedInDescendingOrder(bullionId, page, limit)
}

func (service *feedsService) DeleteById(id string, bullionId string, adminId string) (*interfaces.FeedsEntity, error) {
	entity, err := service.feedsRepo.FindOne(id)
	if err != nil {
		return entity, err
	}
	if entity.BullionId != bullionId {
		return nil, &interfaces.RequestError{
			StatusCode: 403,
			Code:       interfaces.ERROR_MISMATCH_BULLION_ID,
			Message:    "You do not have access to this Feed",
			Name:       "ERROR_MISMATCH_BULLION_ID",
		}
	}
	err = service.feedsRepo.DeleteById(entity.ID)
	if err != nil {
		return entity, err
	}
	event := events.CreateDeleteFeedEvent(entity.FeedsBase, entity.ID, adminId)
	service.eventBus.Publish(event)
	return entity, err
}
