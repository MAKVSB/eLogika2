package wrappers

import (
	"errors"

	"elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/auth/tokens"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"github.com/gin-gonic/gin"
)

type HandlerWithUserDataRole func(c *gin.Context, userData dtos.LoggedUserDTO, role enums.CourseUserRoleEnum) *common.ErrorResponse

func StringToRoleEnum(str string) (enums.CourseUserRoleEnum, error) {
	switch str {
	case string(enums.CourseUserRoleStudent):
		return enums.CourseUserRoleStudent, nil
	case string(enums.CourseUserRoleTutor):
		return enums.CourseUserRoleTutor, nil
	case string(enums.CourseUserRoleGarant):
		return enums.CourseUserRoleGarant, nil
	// case string(enums.CourseUserRoleSecretary):
	// 	return enums.CourseUserRoleSecretary, nil
	case string(enums.CourseUserRoleAdmin):
		return enums.CourseUserRoleAdmin, nil
	default:
		return enums.CourseUserRoleStudent, errors.New("role does not exists")
	}
}

func WithUserDataRole(handler HandlerWithUserDataRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("user")
		if !ok {
			c.AbortWithStatusJSON(403, common.ErrorResponse{
				Code:    403,
				Message: "Access denied",
			})
			return
		}

		roleHeaderEnum, err := StringToRoleEnum(c.GetHeader("X-AS-ROLE"))
		if err != nil {
			c.AbortWithStatusJSON(403, common.ErrorResponse{
				Code:    403,
				Message: "Role does not exists",
			})
			return
		}

		switch v := val.(type) {
		case *tokens.AccessToken:
			res := handler(c, v.UserData, roleHeaderEnum)
			if res != nil {
				c.AbortWithStatusJSON(res.Code, res)
				return
			}
		case *tokens.ApiToken:
			res := handler(c, v.UserData, roleHeaderEnum)
			if res != nil {
				c.AbortWithStatusJSON(res.Code, res)
				return
			}
		default:
			c.AbortWithStatusJSON(403, common.ErrorResponse{
				Message: "Invalid user claims",
			})
			return
		}
	}
}
