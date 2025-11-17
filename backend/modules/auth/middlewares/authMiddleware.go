package middlewares

import (
	"strings"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/auth/tokens"
	"elogika.vsb.cz/backend/modules/common"
	"github.com/gin-gonic/gin"
)

func ModifyLogEntry(c *gin.Context, token models.AuthToken) bool {
	logEntry := models.GetAccesLogEntry(c)
	if logEntry == nil {
		return false
	}
	logEntry.UserID = &token.UserID
	logEntry.TokenID = &token.TokenID
	return true
}

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
				c.AbortWithStatusJSON(err.Code, err)
				return
			}

			if accessToken.IsRevoked() {
				c.AbortWithStatusJSON(401, common.ErrorResponse{
					Code:    401,
					Message: "Token revoked",
				})
				return
			}

			ok := ModifyLogEntry(c, accessToken.AuthToken)
			if !ok {
				c.AbortWithStatusJSON(401, common.ErrorResponse{
					Code:    401,
					Message: "Failed to assign user data to access log",
				})
				return
			}
			c.Set("user", accessToken)
			c.Next()
		} else if strings.HasPrefix(authHeader, "Bearer api_") {
			apiToken := &tokens.ApiToken{}
			err := apiToken.Parse(strings.TrimPrefix(authHeader, "Bearer api_"), false)
			if err != nil {
				c.AbortWithStatusJSON(err.Code, err)
				return
			}

			if apiToken.IsRevoked() {
				c.AbortWithStatusJSON(401, common.ErrorResponse{
					Code:    401,
					Message: "Token revoked",
				})
				return
			}

			ok := ModifyLogEntry(c, apiToken.AuthToken)
			if !ok {
				c.AbortWithStatusJSON(401, common.ErrorResponse{
					Code:    401,
					Message: "Failed to assign user data to access log",
				})
				return
			}
			c.Set("user", apiToken)
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, common.ErrorResponse{
				Code:    401,
				Message: "Invalid token format",
			})
			return
		}
	}
}
