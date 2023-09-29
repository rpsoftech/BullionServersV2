package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility/jwt"
)

// fiber middleware for jwt
func TokenDecrypter(c *fiber.Ctx) (err error) {
	reqHeaders := c.GetReqHeaders()
	tokenString, foundToken := reqHeaders[env.RequestTokenHeaderKey]
	if !foundToken {
		return c.Next()
	}
	if tokenString == "" {
		return c.Next()
	}
	userRolesCustomClaim, errs := services.AccessTokenService.VerifyToken(tokenString)
	if errs != nil {
		if errs == jwt.ErrInvalidSignatureMethod {
			err = &interfaces.RequestError{
				StatusCode: 401,
				Code:       interfaces.ERROR_INVALID_TOKEN_SIGNATURE,
				Message:    "Invalid Token Signature",
				Name:       "JwtInvalidTokenSignature",
			}
		}
		if errs == jwt.ErrTokenExpired {
			err = &interfaces.RequestError{
				StatusCode: 401,
				Code:       interfaces.ERROR_TOKEN_EXPIRED,
				Message:    "Invalid Token Expired",
				Name:       "JwtInvalidTokenExpired",
			}
		}
		if err == nil {
			err = &interfaces.RequestError{
				StatusCode: 401,
				Code:       interfaces.ERROR_INVALID_TOKEN_SIGNATURE,
				Message:    "Invalid Token Body",
				Name:       "JwtInvalidTokenBody",
			}
		}
		return err
	}
	mappedClaim, ok := userRolesCustomClaim.Claims.(map[string]interface{})
	if !ok {
		err = &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_INVALID_TOKEN_SIGNATURE,
			Message:    "Invalid Token Body",
			Name:       "JwtInvalidTokenBody",
		}
		return err
	}
	role, ok := mappedClaim["role"].(string)

	if !ok || !interfaces.ValidateEnumUserRole(role) {
		err = &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_TOKEN_ROLE_NOT_FOUND,
			Message:    "Invalid Token Role Or Not Found",
			Name:       "INVALID_TOKEN_ROLE",
		}

		return err
	}

	c.Locals(interfaces.REQ_LOCAL_KEY_ROLE, role)
	c.Locals(interfaces.REQ_LOCAL_KEY_TOKEN_RAW_DATA, mappedClaim)
	return c.Next()
	// TODO: Base on role decrypt interface of users

}
