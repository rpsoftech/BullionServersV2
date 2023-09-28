package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
)

// fiber middleware for jwt
func TokenDecrypter(c *fiber.Ctx) (err error) {
	reqHeaders := c.GetReqHeaders()
	tokenString, foundToken := reqHeaders[env.RequestTokenHeaderKey]
	if !foundToken {
		c.Next()
		return
	}
	if tokenString == "" {
		c.Next()
		return
	}
	userRolesCustomClaim, err := services.AccessTokenService.VerifyToken(tokenString)
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

	// TODO: Base on role decrypt interface of users

	return
}

// func JwtAuthMiddleware(isAdmin bool) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenString := strings.Split(c.GetHeader("Authorization"), " ")[1]
// 		if tokenString == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 			return
// 		}

// 		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
// 			// Check the signing method
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 			}
// 			// Set the secret key for verification
// 			return global.JWTKEY, nil
// 		})

// 		if err != nil || !token.Valid {

// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": err.Error()})
// 			return
// 		}

// 		if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
// 			if isAdmin && claims.IsAdmin() == false {
// 				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token Value"})
// 				return
// 			}
// 			c.Set("User", claims)
// 		} else {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token Value"})
// 			return
// 		}
// 		// Token is valid
// 		c.Next()
// 	}
// }
