package interfaces

import "github.com/golang-jwt/jwt/v5"

type GeneralUserAccessRefreshToken struct {
	*jwt.RegisteredClaims
	Role      UserRoles `json:"role"`
	UserId    string    `json:"userId"`
	BullionId string    `json:"bullionId"`
}
