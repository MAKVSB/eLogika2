package wrappers

import (
	"elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/auth/tokens"
	"elogika.vsb.cz/backend/modules/common"
	"github.com/gin-gonic/gin"
)

type HandlerWithUserDataAndToken func(c *gin.Context, userData dtos.LoggedUserDTO, token tokens.Token)

func WithUserDataAndToken(handler HandlerWithUserDataAndToken) gin.HandlerFunc {
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
			handler(c, v.UserData, v)
		case *tokens.ApiToken:
			handler(c, v.UserData, v)
		default:
			c.AbortWithStatusJSON(403, common.ErrorResponse{
				Message: "Invalid user claims",
			})
			return
		}
	}
}
