package middleware

import (
	"fmt"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
)

type roleCheckerMiddlewareWithRolesArray struct {
	roles []string
}

var (
	allUserRoles = []string{
		interfaces.ROLE_RATE_ADMIN.String(),
		interfaces.ROLE_SUPER_ADMIN.String(),
		interfaces.ROLE_ADMIN.String(),
		interfaces.ROLE_GENERAL_USER.String(),
		interfaces.ROLE_TRADE_USER.String(),
	}

	AllowAllUsers = roleCheckerMiddlewareWithRolesArray{
		roles: allUserRoles,
		// Middleware: roleCheckerAndValidator(allUserRoles),
	}
)

func (cc *roleCheckerMiddlewareWithRolesArray) Validate(c *fiber.Ctx) (err error) {
	roleFromLocal := c.Locals(interfaces.REQ_LOCAL_KEY_ROLE)
	if roleFromLocal == nil {
		err = &interfaces.RequestError{
			StatusCode: 403,
			Code:       interfaces.ERROR_TOKEN_ROLE_NOT_FOUND,
			Message:    "Invalid Token Role Or Not Found",
			Name:       "INVALID_TOKEN_ROLE",
		}
		return
	}
	role, ok := roleFromLocal.(string)
	if !ok {
		err = &interfaces.RequestError{
			StatusCode: 403,
			Code:       interfaces.ERROR_ROLE_NOT_EXISTS,
			Message:    "Token Role Should be string",
			Name:       "INVALID_TOKEN_ROLE_FORMAT",
		}
		return
	}
	if role == string(interfaces.ROLE_GOD) {
		return c.Next()
	}
	if !slices.Contains(cc.roles, role) {
		err = &interfaces.RequestError{
			StatusCode: 403,
			Code:       interfaces.ERROR_ROLE_NOT_AUTHORIZED,
			Message:    fmt.Sprintf("%s is not allowed for this route", role),
			Name:       "INVALID_TOKEN_ROLE_FORMAT",
		}
		return
	}
	return c.Next()
}
