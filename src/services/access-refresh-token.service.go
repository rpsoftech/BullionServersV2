package services

import (
	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/utility/jwt"
)

var AccessTokenService *jwt.TokenService
var RefreshTokenService *jwt.TokenService

func init() {
	AccessTokenService = &jwt.TokenService{
		SigningKey: []byte(env.Env.ACCESS_TOKEN_KEY),
	}
	RefreshTokenService = &jwt.TokenService{
		SigningKey: []byte(env.Env.REFRESH_TOKEN_KEY),
	}
}
