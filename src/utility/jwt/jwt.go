package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rpsoftech/bullion-server/src/interfaces"
)

var (
	ErrTokenExpired           = jwt.ErrTokenExpired
	ErrInvalidToken           = errors.New("invalid token")
	ErrJwt                    = errors.New("something went wrong with Token")
	ErrSignatureInvalid       = jwt.ErrSignatureInvalid
	ErrInvalidSignatureMethod = errors.New("unexpected signing method")
)

type GeneralUserAccessRefreshToken struct {
	*jwt.RegisteredClaims
	UserId     string                 `json:"userId" validate:"required"`
	BullionId  string                 `json:"bullionId" validate:"required"`
	Role       interfaces.UserRoles   `json:"role" validate:"required"`
	ExtraClaim map[string]interface{} `json:"extraClaim,omitempty"`
}
type TokenService struct {
	SigningKey []byte
}

func (t *TokenService) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.SigningKey)
}
func (t *TokenService) VerifyToken(token string) (*GeneralUserAccessRefreshToken, error) {
	claimRaw, err := jwt.ParseWithClaims(token, &GeneralUserAccessRefreshToken{}, func(token *jwt.Token) (any, error) {
		// return t.SigningKey, nil
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &interfaces.RequestError{
				StatusCode: 401,
				Code:       interfaces.ERROR_INVALID_TOKEN_SIGNATURE,
				Message:    "Invalid Signature Method",
				Name:       "ERROR_INVALID_TOKEN_SIGNATURE",
				Extra:      ErrInvalidSignatureMethod,
			}
		}
		// Set the secret key for verification
		return t.SigningKey, nil
	})

	fmt.Sprintf("claimRaw: %v", claimRaw.Claims)

	claim, ok := claimRaw.Claims.(*GeneralUserAccessRefreshToken)
	if (!claimRaw.Valid || !ok) && err != nil {
		err = &interfaces.RequestError{
			StatusCode: 401,
			Code:       interfaces.ERROR_INVALID_TOKEN_SIGNATURE,
			Message:    "Error InValid Token",
			Name:       "ERROR_INVALID_TOKEN",
			Extra:      err,
		}
	}
	if err == ErrTokenExpired {
		err = &interfaces.RequestError{
			StatusCode: 403,
			Code:       interfaces.ERROR_TOKEN_EXPIRED,
			Message:    "Refresh Token Expired",
			Name:       "ERROR_TOKEN_EXPIRED",
			Extra:      err,
		}
	}
	return claim, err
}
