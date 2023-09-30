package services

import "github.com/rpsoftech/bullion-server/src/mongodb/repos"

type generalUserService struct {
	GeneralUserRepo     *repos.GeneralUserRepoStruct
	BullionSiteInfoRepo *repos.BullionSiteInfoRepoStruct
}

var GeneralUserService *generalUserService

func init() {
	GeneralUserService = &generalUserService{
		GeneralUserRepo:     repos.GeneralUserRepo,
		BullionSiteInfoRepo: repos.BullionSiteInfoRepo,
	}
}

func (service *generalUserService) RegisterNew(bullionId string) {

}
