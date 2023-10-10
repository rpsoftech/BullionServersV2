package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
	"github.com/rpsoftech/bullion-server/src/utility"
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
	userRolesCustomClaim, localErr := services.AccessTokenService.VerifyToken(tokenString)
	if localErr != nil {
		c.Locals(interfaces.REQ_LOCAL_ERROR_KEY, localErr)
		return c.Next()
	}
	if localErr := utility.ValidateStructAndReturnReqError(&userRolesCustomClaim, &interfaces.RequestError{
		StatusCode: 401,
		Code:       interfaces.ERROR_INVALID_TOKEN_SIGNATURE,
		Message:    "Invalid Token Structure",
		Name:       "ERROR_INVALID_TOKEN_SIGNATURE",
	}); localErr != nil {
		c.Locals(interfaces.REQ_LOCAL_ERROR_KEY, localErr)
		return c.Next()
	}
	role := userRolesCustomClaim.Role.String()
	if !interfaces.ValidateEnumUserRole(role) {
		err = &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_TOKEN_ROLE_NOT_FOUND,
			Message:    "Invalid Token Role Or Not Found",
			Name:       "INVALID_TOKEN_ROLE",
		}

		c.Locals(interfaces.REQ_LOCAL_ERROR_KEY, err)
		return c.Next()
	}

	c.Locals(interfaces.REQ_LOCAL_KEY_ROLE, role)
	c.Locals(interfaces.REQ_LOCAL_KEY_TOKEN_RAW_DATA, userRolesCustomClaim)
	return c.Next()
	// TODO: Base on role decrypt interface of users
}

func AllowOnlyValidTokenMiddleWare(c *fiber.Ctx) error {
	jwtRawFromLocal := c.Locals(interfaces.REQ_LOCAL_KEY_TOKEN_RAW_DATA)
	localError := c.Locals(interfaces.REQ_LOCAL_ERROR_KEY)
	if jwtRawFromLocal == nil {
		if localError != nil {
			err, ok := localError.(*interfaces.RequestError)
			if !ok {
				return &interfaces.RequestError{
					StatusCode: 403,
					Code:       interfaces.ERROR_TOKEN_EXPIRED,
					Message:    "Cannot Cast Error Token",
					Name:       "NOT_VALID_DECRYPTED_TOKEN",
				}
			}
			return err
		} else {

			err := &interfaces.RequestError{
				StatusCode: 403,
				Code:       interfaces.ERROR_TOKEN_EXPIRED,
				Message:    "Invalid Token or token expired",
				Name:       "NOT_VALID_DECRYPTED_TOKEN",
			}
			return err
		}
	}
	return c.Next()
}
