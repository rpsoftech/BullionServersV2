package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mitchellh/mapstructure"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
	localJwt "github.com/rpsoftech/bullion-server/src/utility/jwt"
	"github.com/rpsoftech/bullion-server/src/validator"
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
	var baseGeneralUser interfaces.GeneralUser
	var entity interfaces.GeneralUserEntity

	Bullion, err := service.BullionSiteInfoRepo.FindOne(bullionId)
	if err != nil {
		return &entity, err
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
	} else {
		baseGeneralUser = interfaces.GeneralUser{
			IsAuto: false,
		}
	}
	baseGeneralUser.RandomPass = faker.Password()
	err = mapstructure.Decode(user, &baseGeneralUser)
	if err != nil {
		return &entity, err
	}
	errs := validator.Validator.Validate(&baseGeneralUser)
	if len(errs) > 0 {
		reqErr := &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "",
			Name:       "INVALID_INPUT",
			Extra:      errs,
		}
		return &entity, reqErr.AppendValidationErrors(errs)
	}
	entity = interfaces.GeneralUserEntity{
		BaseEntity:  interfaces.BaseEntity{},
		GeneralUser: baseGeneralUser,
		UserRolesInterface: interfaces.UserRolesInterface{
			Role: interfaces.ROLE_GENERAL_USER,
		},
	}
	entity.CreateNewId()

	errs = validator.Validator.Validate(&entity)
	if len(errs) > 0 {
		reqErr := &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "",
			Name:       "INVALID_INPUT",
			Extra:      errs,
		}
		return &entity, reqErr.AppendValidationErrors(errs)
	}
	service.GeneralUserRepo.Save(&entity)
	_, err = service.sendApprovalRequest(&entity, Bullion)
	if err != nil {
		return &entity, err
	}
	return &entity, err
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
	reqEntity = &interfaces.GeneralUserReqEntity{
		GeneralUserId: user.ID,
		BullionId:     bullion.ID,
		Status:        interfaces.GENERAL_USER_AUTH_STATUS_REQUESTED,
	}
	if bullion.GeneralUserInfo.AutoApprove {
		reqEntity.Status = interfaces.GENERAL_USER_AUTH_STATUS_AUTHORIZED
	}
	reqEntity.CreateNewId()
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
	var tokenResponse *interfaces.TokenResponseBody
	userEntity, err := service.GetGeneralUserDetailsByIdPassword(userId, password)
	if err != nil {
		return tokenResponse, err
	}
	return service.validateApprovalAndGenerateTokenStage2(userEntity, bullionId)
}
func (service *generalUserService) validateApprovalAndGenerateTokenStage2(user *interfaces.GeneralUserEntity, bullionId string) (*interfaces.TokenResponseBody, error) {
	var tokenResponse *interfaces.TokenResponseBody
	reqEntity, err := service.generalUserReqRepo.FindOneByGeneralUserIdAndBullionId(user.ID, bullionId)
	if err != nil || reqEntity == nil {
		err = &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_GENERAL_USER_REQ_NOT_FOUND,
			Message:    "REQUEST DOES NOT EXISTS",
			Name:       "ERROR_GENERAL_USER_REQ_NOT_FOUND",
		}
		return tokenResponse, err
	}
	if reqEntity.Status == interfaces.GENERAL_USER_AUTH_STATUS_REQUESTED {
		err = &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_GENERAL_USER_REQ_PENDING,
			Message:    "REQUEST PENDING",
			Name:       "ERROR_GENERAL_USER_REQ_PENDING",
		}
		return tokenResponse, err
	}
	if reqEntity.Status == interfaces.GENERAL_USER_AUTH_STATUS_REJECTED {
		err = &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_GENERAL_USER_REQ_REJECTED,
			Message:    "REQUEST REJECTED",
			Name:       "ERROR_GENERAL_USER_REQ_PENDING",
		}
		return tokenResponse, err
	}

	return service.generateTokens(user.ID, bullionId)
}

func (service *generalUserService) generateTokens(userId string, bullionId string) (*interfaces.TokenResponseBody, error) {
	var tokenResponse *interfaces.TokenResponseBody
	now := time.Now()
	accessToken, err := AccessTokenService.GenerateToken(localJwt.GeneralUserAccessRefreshToken{
		UserId:    userId,
		BullionId: bullionId,
		Role:      interfaces.ROLE_GENERAL_USER,
		RegisteredClaims: &jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: now},
			ExpiresAt: &jwt.NumericDate{Time: now.Add(time.Minute * 30)},
		},
	})
	if err != nil {
		err = &interfaces.RequestError{
			Code:    interfaces.ERROR_INTERNAL_SERVER,
			Message: "JWT ACCESS TOKEN GENERATION ERROR",
			Name:    "ERROR_INTERNAL_ERROR",
			Extra:   err,
		}
		return tokenResponse, err
	}
	refreshToken, err := RefreshTokenService.GenerateToken(localJwt.GeneralUserAccessRefreshToken{
		UserId:    userId,
		BullionId: bullionId,
		Role:      interfaces.ROLE_GENERAL_USER,
		RegisteredClaims: &jwt.RegisteredClaims{
			IssuedAt: &jwt.NumericDate{Time: now},
			// ExpiresAt: &jwt.NumericDate{Time: now.Add(time.Hour * 24 * 30)},
		},
	})
	if err != nil {
		err = &interfaces.RequestError{
			Code:    interfaces.ERROR_INTERNAL_SERVER,
			Message: "JWT ACCESS TOKEN GENERATION ERROR",
			Name:    "ERROR_INTERNAL_ERROR",
			Extra:   err,
		}
		return tokenResponse, err
	}
	firebaseToken, err := FirebaseAuthService.GenerateCustomToken(userId, map[string]interface{}{
		"userId":    userId,
		"bullionId": bullionId,
		"role":      interfaces.ROLE_GENERAL_USER,
	})
	tokenResponse = &interfaces.TokenResponseBody{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		FirebaseToken: firebaseToken,
	}
	return tokenResponse, err
}
func (service *generalUserService) RefreshToken(token string) (*interfaces.TokenResponseBody, error) {
	var tokenResponse *interfaces.TokenResponseBody
	_, err := RefreshTokenService.VerifyToken(token)

	// d, _ := json.Marshal(rr)
	// fmt.Printf("%+v\n", rr)
	// fmt.Printf("%s\n", d)
	return tokenResponse, err
}
