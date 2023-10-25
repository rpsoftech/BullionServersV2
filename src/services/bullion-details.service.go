package services

import (
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type bullionDetailsService struct {
	BullionSiteInfoRepo *repos.BullionSiteInfoRepoStruct
}

var BullionDetailsService *bullionDetailsService

func init() {
	BullionDetailsService = &bullionDetailsService{
		BullionSiteInfoRepo: repos.BullionSiteInfoRepo,
	}
	println("Bullion Site Details Initialized")
}

func (service *bullionDetailsService) GetBullionDetailsByShortName(shortName string) (*interfaces.BullionSiteInfoEntity, error) {
	bullion, err := service.BullionSiteInfoRepo.FindByShortName(shortName)
	if err != nil {
		return nil, err
	}
	return bullion, nil
}
func (service *bullionDetailsService) GetBullionDetailsByBullionId(id string) (*interfaces.BullionSiteInfoEntity, error) {
	bullion, err := service.BullionSiteInfoRepo.FindOne(id)
	if err != nil {
		return nil, err
	}
	return bullion, nil
}
