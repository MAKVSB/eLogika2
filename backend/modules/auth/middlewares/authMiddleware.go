package middlewares

import (
	"strings"

	"elogika.vsb.cz/backend/modules/auth/tokens"
	"elogika.vsb.cz/backend/modules/common"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(500, gin.H{
				"code":    500,
				"message": "Missing token header",
			})
			return
		}

		if strings.HasPrefix(authHeader, "Bearer at_") {
			accessToken := &tokens.AccessToken{}
			err := accessToken.Parse(strings.TrimPrefix(authHeader, "Bearer at_"), false)
			if err != nil {
				c.AbortWithStatusJSON(401, err)
				return
			}

			if accessToken.IsRevoked() {
				c.AbortWithStatusJSON(401, common.ErrorResponse{
					Message: "Token revoked",
				})
				return
			}

			c.Set("user", accessToken)
			c.Next()
		} else if strings.HasPrefix(authHeader, "Bearer api_") {
			apiToken := &tokens.ApiToken{}
			err := apiToken.Parse(strings.TrimPrefix(authHeader, "Bearer api_"), false)
			if err != nil {
				c.AbortWithStatusJSON(401, err)
				return
			}

			if apiToken.IsRevoked() {
				c.AbortWithStatusJSON(401, common.ErrorResponse{
					Message: "Token revoked",
				})
				return
			}

			c.Set("user", apiToken)
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": "Invalid token format",
			})
			return
		}
	}
}
