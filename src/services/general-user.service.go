package services

import (
	"fmt"
	"math/rand"

	"github.com/go-faker/faker/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
	"github.com/rpsoftech/bullion-server/src/utility"
)

type generalUserService struct {
	generalUserReqRepo  *repos.GeneralUserReqRepoStruct
	GeneralUserRepo     *repos.GeneralUserRepoStruct
	BullionSiteInfoRepo *repos.BullionSiteInfoRepoStruct
}

var GeneralUserService *generalUserService

func init() {
	GeneralUserService = &generalUserService{
		GeneralUserRepo:     repos.GeneralUserRepo,
		BullionSiteInfoRepo: repos.BullionSiteInfoRepo,
		generalUserReqRepo:  repos.GeneralUserReqRepo,
	}
}

func (service *generalUserService) RegisterNew(bullionId string, user interface{}) (*interfaces.GeneralUserEntity, error) {
	Bullion, err := service.BullionSiteInfoRepo.FindOne(bullionId)
	if err != nil {
		return nil, err
	}

	baseGeneralUser := interfaces.GeneralUser{
		IsAuto: false,
	}
	if Bullion.GeneralUserInfo.AutoLogin {
		baseGeneralUser = interfaces.GeneralUser{
			FirstName:     faker.FirstName(),
			LastName:      faker.LastName(),
			FirmName:      faker.Username(),
			ContactNumber: faker.Phonenumber(),
			GstNumber:     fmt.Sprintf("%dAAAAA%dA1ZA", rand.Intn(99-10)+10, rand.Intn(9999-1000)+1000),
			OS:            "AUTO",
			IsAuto:        true,
			DeviceId:      faker.UUIDDigit(),
			DeviceType:    interfaces.DEVICE_TYPE_IOS,
		}
	}

	baseGeneralUser.RandomPass = faker.Password()
	err = mapstructure.Decode(user, &baseGeneralUser)
	if err != nil {
		return nil, err
	}

	if err := utility.ValidateReqInput(&baseGeneralUser); err != nil {
		return nil, err
	}

	entity := interfaces.CreateNewGeneralUser(baseGeneralUser)
	if err := utility.ValidateReqInput(&entity); err != nil {
		return nil, err
	}

	service.GeneralUserRepo.Save(entity)
	_, err = service.sendApprovalRequest(entity, Bullion)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (service *generalUserService) CreateApprovalRequest(userId string, password string, bullionId string) (reqEntity *interfaces.GeneralUserReqEntity, err error) {
	var userEntity *interfaces.GeneralUserEntity
	var bullionEntity *interfaces.BullionSiteInfoEntity
	if userEntity, err = service.GetGeneralUserDetailsByIdPassword(userId, password); err == nil {
		if bullionEntity, err = service.BullionSiteInfoRepo.FindOne(bullionId); err == nil {
			reqEntity, err = service.sendApprovalRequest(userEntity, bullionEntity)
		}
	}
	return
}
func (service *generalUserService) sendApprovalRequest(user *interfaces.GeneralUserEntity, bullion *interfaces.BullionSiteInfoEntity) (reqEntity *interfaces.GeneralUserReqEntity, err error) {
	reqEntity, err = service.generalUserReqRepo.FindOneByGeneralUserIdAndBullionId(user.ID, bullion.ID)
	if err == nil {
		err = &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_GENERAL_USER_REQ_EXISTS,
			Message:    "REQUEST ALREADY EXISTS",
			Name:       "ERROR_GENERAL_USER_REQ_EXISTS",
		}
		return
	} else {
		err = nil
	}
	reqEntity = interfaces.CreateNewGeneralUserReq(user.ID, bullion.ID, interfaces.GENERAL_USER_AUTH_STATUS_REQUESTED)
	if bullion.GeneralUserInfo.AutoApprove {
		reqEntity.Status = interfaces.GENERAL_USER_AUTH_STATUS_AUTHORIZED
	}
	reqEntity, err = service.generalUserReqRepo.Save(reqEntity)
	return
}
func (service *generalUserService) GetGeneralUserDetailsByIdPassword(id string, password string) (*interfaces.GeneralUserEntity, error) {
	entity, err := service.GeneralUserRepo.FindOne(id)
	if err != nil {
		return entity, err
	}
	if entity.RandomPass != password {
		err = &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_GENERAL_USER_INVALID_PASSWORD,
			Message:    fmt.Sprintf("GeneralUser Entity invalid password %s ", password),
			Name:       "ERROR_GENERAL_USER_INVALID_PASSWORD",
		}
		return entity, err
	}
	return entity, err
}

func (service *generalUserService) ValidateApprovalAndGenerateToken(userId string, password string, bullionId string) (*interfaces.TokenResponseBody, error) {
	userEntity, err := service.GetGeneralUserDetailsByIdPassword(userId, password)
	if err != nil {
		return nil, err
	}
	return service.validateApprovalAndGenerateTokenStage2(userEntity, bullionId)
}
func (service *generalUserService) validateApprovalAndGenerateTokenStage2(user *interfaces.GeneralUserEntity, bullionId string) (*interfaces.TokenResponseBody, error) {
	reqEntity, err := service.generalUserReqRepo.FindOneByGeneralUserIdAndBullionId(user.ID, bullionId)
	if err != nil || reqEntity == nil {
		return nil, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_GENERAL_USER_REQ_NOT_FOUND,
			Message:    "REQUEST DOES NOT EXISTS",
			Name:       "ERROR_GENERAL_USER_REQ_NOT_FOUND",
		}
	}

	switch reqEntity.Status {
	case interfaces.GENERAL_USER_AUTH_STATUS_REQUESTED:
		return nil, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_GENERAL_USER_REQ_PENDING,
			Message:    "REQUEST PENDING",
			Name:       "ERROR_GENERAL_USER_REQ_PENDING",
		}
	case interfaces.GENERAL_USER_AUTH_STATUS_REJECTED:
		return nil, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_GENERAL_USER_REQ_REJECTED,
			Message:    "REQUEST REJECTED",
			Name:       "ERROR_GENERAL_USER_REQ_REJECTED",
		}
	case interfaces.GENERAL_USER_AUTH_STATUS_AUTHORIZED:
		return service.generateTokens(user.ID, bullionId)
	default:
		return nil, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_GENERAL_USER_INVALID_STATUS,
			Message:    "Invalid Request Status",
			Name:       "ERROR_GENERAL_USER_INVALID_STATUS",
		}
	}
}

func (service *generalUserService) generateTokens(userId string, bullionId string) (*interfaces.TokenResponseBody, error) {
	return generateTokens(userId, bullionId, interfaces.ROLE_GENERAL_USER)
}
func (service *generalUserService) RefreshToken(token string) (*interfaces.TokenResponseBody, error) {
	var tokenResponse *interfaces.TokenResponseBody
	tokenBody, err := RefreshTokenService.VerifyToken(token)
	if err != nil {
		return tokenResponse, err
	}

	generalUserEntity, err := service.GeneralUserRepo.FindOne(tokenBody.UserId)
	if err != nil {
		return tokenResponse, err
	}
	return service.validateApprovalAndGenerateTokenStage2(generalUserEntity, tokenBody.BullionId)
}
