package wrappers

import (
	"elogika.vsb.cz/backend/modules/auth/tokens"
	"elogika.vsb.cz/backend/modules/common"
	"github.com/gin-gonic/gin"
)

type HandlerWithToken func(c *gin.Context, token tokens.Token)

func WithToken(handler HandlerWithToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("user")
		if !ok {
			c.AbortWithStatusJSON(403, common.ErrorResponse{
				Message: "Access denied",
			})
			return
		}

		switch v := val.(type) {
		case *tokens.AccessToken:
			handler(c, v)
		case *tokens.ApiToken:
			handler(c, v)
		default:
			c.AbortWithStatusJSON(403, common.ErrorResponse{
				Message: "Invalid user claims",
			})
			return
		}
	}
}
