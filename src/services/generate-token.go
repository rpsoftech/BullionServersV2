package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	localJwt "github.com/rpsoftech/bullion-server/src/utility/jwt"
)

func generateTokens(userId string, bullionId string, role interfaces.UserRoles) (*interfaces.TokenResponseBody, error) {
	var tokenResponse *interfaces.TokenResponseBody
	now := time.Now()
	accessToken, err := AccessTokenService.GenerateToken(localJwt.GeneralUserAccessRefreshToken{
		UserId:    userId,
		BullionId: bullionId,
		Role:      role,
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
		Role:      role,
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
		"role":      role,
	})
	tokenResponse = &interfaces.TokenResponseBody{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		FirebaseToken: firebaseToken,
	}
	return tokenResponse, err
}
