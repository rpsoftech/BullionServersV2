package services

import "github.com/rpsoftech/bullion-server/src/mongodb/repos"

type adminUserService struct {
	adminUserRepo       *repos.GeneralUserReqRepoStruct
	GeneralUserRepo     *repos.GeneralUserRepoStruct
	BullionSiteInfoRepo *repos.BullionSiteInfoRepoStruct
}

var AdminUserService *adminUserService
