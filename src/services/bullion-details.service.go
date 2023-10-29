package services

import (
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type bullionDetailsService struct {
	BullionSiteInfoRepo           *repos.BullionSiteInfoRepoStruct
	billionSiteInfoMapById        map[string]*interfaces.BullionSiteInfoEntity
	billionSiteInfoMapByShortName map[string]*interfaces.BullionSiteInfoEntity
}

var BullionDetailsService *bullionDetailsService

func init() {
	BullionDetailsService = &bullionDetailsService{
		BullionSiteInfoRepo:    repos.BullionSiteInfoRepo,
		billionSiteInfoMapById: make(map[string]*interfaces.BullionSiteInfoEntity),
	}
	println("Bullion Site Details Initialized")
}

func (service *bullionDetailsService) GetBullionDetailsByShortName(shortName string) (*interfaces.BullionSiteInfoEntity, error) {
	if bullion, ok := service.billionSiteInfoMapByShortName[shortName]; ok {
		return bullion, nil
	}
	bullion, err := service.BullionSiteInfoRepo.FindByShortName(shortName)
	if err != nil {
		return nil, err
	}
	service.billionSiteInfoMapById[shortName] = bullion
	return bullion, nil
}
func (service *bullionDetailsService) GetBullionDetailsByBullionId(id string) (*interfaces.BullionSiteInfoEntity, error) {
	if bullion, ok := service.billionSiteInfoMapById[id]; ok {
		return bullion, nil
	}
	bullion, err := service.BullionSiteInfoRepo.FindOne(id)
	if err != nil {
		return nil, err
	}
	service.billionSiteInfoMapById[id] = bullion
	return bullion, nil
}
