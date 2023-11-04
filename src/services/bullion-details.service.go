package services

import (
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type bullionDetailsService struct {
	BullionSiteInfoRepo           *repos.BullionSiteInfoRepoStruct
	bullionSiteInfoMapById        map[string]*interfaces.BullionSiteInfoEntity
	bullionSiteInfoMapByShortName map[string]*interfaces.BullionSiteInfoEntity
}

var BullionDetailsService *bullionDetailsService

func init() {
	getBullionService()
}

func getBullionService() *bullionDetailsService {
	if BullionDetailsService == nil {
		BullionDetailsService = &bullionDetailsService{
			BullionSiteInfoRepo:           repos.BullionSiteInfoRepo,
			bullionSiteInfoMapById:        make(map[string]*interfaces.BullionSiteInfoEntity),
			bullionSiteInfoMapByShortName: make(map[string]*interfaces.BullionSiteInfoEntity),
		}
		println("Bullion Site Details Initialized")
	}
	return BullionDetailsService
}

func (service *bullionDetailsService) GetBullionDetailsByShortName(shortName string) (*interfaces.BullionSiteInfoEntity, error) {
	if bullion, ok := service.bullionSiteInfoMapByShortName[shortName]; ok {
		return bullion, nil
	}
	bullion, err := service.BullionSiteInfoRepo.FindByShortName(shortName)
	if err != nil {
		return nil, err
	}
	service.bullionSiteInfoMapById[shortName] = bullion
	return bullion, nil
}
func (service *bullionDetailsService) GetBullionDetailsByBullionId(id string) (*interfaces.BullionSiteInfoEntity, error) {
	if bullion, ok := service.bullionSiteInfoMapById[id]; ok {
		return bullion, nil
	}
	bullion, err := service.BullionSiteInfoRepo.FindOne(id)
	if err != nil {
		return nil, err
	}
	service.bullionSiteInfoMapById[id] = bullion
	return bullion, nil
}

func (service *bullionDetailsService) UpdateBullionSiteDetails(details *interfaces.BullionSiteInfoEntity) (*interfaces.BullionSiteInfoEntity, error) {
	service.bullionSiteInfoMapById[details.ID] = details
	service.bullionSiteInfoMapByShortName[details.ShortName] = details
	return service.BullionSiteInfoRepo.Save(details)
}
