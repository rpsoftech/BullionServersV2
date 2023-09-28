package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired           = jwt.ErrTokenExpired
	ErrInvalidToken           = errors.New("invalid token")
	ErrJwt                    = errors.New("something went wrong with Token")
	ErrSignatureInvalid       = jwt.ErrSignatureInvalid
	ErrInvalidSignatureMethod = errors.New("unexpected signing method")
)

type CustomClaims struct {
	Claims interface{}
	jwt.RegisteredClaims
}

type TokenService struct {
	SigningKey []byte
}

func (t *TokenService) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.SigningKey)
}
func (t *TokenService) VerifyToken(token string) (claim CustomClaims, errs error) {
	validatedToken, err := jwt.ParseWithClaims(token, &claim, func(token *jwt.Token) (any, error) {
		// return t.SigningKey, nil
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignatureMethod
		}
		// Set the secret key for verification
		return t.SigningKey, nil
	})

	if !validatedToken.Valid && err != nil {
		err = ErrInvalidToken
	}
	errs = err
	return claim, errs
}
