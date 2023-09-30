package services

import (
	"fmt"
	"math/rand"

	"github.com/go-faker/faker/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
	"github.com/rpsoftech/bullion-server/src/validator"
)

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
			RandomPass:    faker.Password(),
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
	// Pending Approval
	return &entity, err
}
