package middlewares

import (
	"elogika.vsb.cz/backend/modules/auth/tokens"
	"elogika.vsb.cz/backend/modules/common"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := &tokens.AccessToken{}
		err := accessToken.Get(c, false)
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
	}
}
