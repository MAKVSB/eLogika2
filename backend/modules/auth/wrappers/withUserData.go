package wrappers

import (
	"net/http"

	"elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/auth/tokens"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"github.com/gin-gonic/gin"
)

type HandlerWithUserData func(c *gin.Context, userData dtos.LoggedUserDTO)

func WithUserData(handler HandlerWithUserData) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("user")
		if !ok {
			c.AbortWithStatusJSON(403, common.ErrorResponse{
				Message: "Access denied",
			})
			return
		}

		roleHeader := c.GetHeader("X-AS-ROLE")

		// You can also check if the header exists:
		if roleHeader != string(enums.CourseUserRoleStudent) &&
			roleHeader != string(enums.CourseUserRoleTutor) &&
			roleHeader != string(enums.CourseUserRoleGarant) &&
			// roleHeader != string(enums.CourseUserRoleSecretary) &&
			roleHeader != string(enums.CourseUserRoleAdmin) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role header not found"})
			return
		}

		switch v := val.(type) {
		case *tokens.AccessToken:
			handler(c, v.UserData)
		case *tokens.ApiToken:
			handler(c, v.UserData)
		default:
			c.AbortWithStatusJSON(403, common.ErrorResponse{
				Message: "Invalid user claims",
			})
			return
		}
	}
}
