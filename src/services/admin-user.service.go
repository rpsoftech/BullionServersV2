package services

import (
	"net/http"

	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type adminUserService struct {
	adminUserRepo *repos.AdminUserRepoStruct
}

var AdminUserService *adminUserService

func init() {
	AdminUserService = &adminUserService{
		adminUserRepo: repos.AdminUserRepo,
	}
}

func (service *adminUserService) ValidateUserAndGenerateToken(uname string, password string, bullionId string) (*interfaces.TokenResponseBody, error) {
	admin, err := service.adminUserRepo.FindOneUserNameAndBullionId(uname, bullionId)
	if err != nil {
		return nil, err
	}
	if !admin.MatchPassword(password) {
		return nil, &interfaces.RequestError{
			StatusCode: http.StatusUnauthorized,
			Code:       interfaces.ERROR_INVALID_PASSWORD,
			Message:    "Invalid Password",
			Name:       "ERROR_INVALID_PASSWORD",
		}
	}
	return generateTokens(admin.ID, bullionId, interfaces.ROLE_ADMIN)
}
