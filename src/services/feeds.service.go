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

func (service *feedsService) AddNewFeedsAndSendNotification(entity *interfaces.FeedsEntity) (*interfaces.FeedsEntity, error) {
	return service.feedsRepo.Save(entity)
}
