package services

import (
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type feedsService struct {
	feedsRepo  *repos.FeedsRepoStruct
	fcmService *firebaseCloudMessagingService
}

var FeedsService *feedsService

func init() {
	FeedsService = &feedsService{
		feedsRepo:  repos.FeedsRepo,
		fcmService: FirebaseCloudMessagingService,
	}
}

func (service *feedsService) SendNotification(bullionId string, feedId string) error {
	entity, err := service.feedsRepo.FindOne(feedId)
	if err != nil {
		return err
	}
	if entity.BullionId != bullionId {
		return &interfaces.RequestError{
			StatusCode: 403,
			Code:       interfaces.ERROR_MISMATCH_BULLION_ID,
			Message:    "Can not send this feed as Notification. You do not have access to this Feed",
			Name:       "ERROR_MISMATCH_BULLION_ID",
		}
	}
	service.fcmService.SendTextNotificationToAll(bullionId, entity.Title, entity.Body, entity.IsHtml)
	return nil
}
func (service *feedsService) AddAndUpdateNewFeeds(entity *interfaces.FeedsEntity) (*interfaces.FeedsEntity, error) {
	return service.feedsRepo.Save(entity)
}
