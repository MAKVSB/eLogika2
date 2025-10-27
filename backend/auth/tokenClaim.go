package auth

import (
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
)

func GetClaimCourseRole(userData dtos.LoggedUserDTO, courseId uint, role enums.CourseUserRoleEnum) *common.ErrorResponse {
	if userData.Type == enums.UserTypeAdmin {
		return nil
	}

	if len(userData.Courses) != 0 {
		for _, item := range userData.Courses {
			if item.ID == courseId {
				if item.HasRole(role) {
					return nil
				} else {
					return &common.ErrorResponse{
						Code:    403,
						Message: "User does not have requested role is course",
					}
				}
			}
		}
		return &common.ErrorResponse{
			Code:    403,
			Message: "User does not belong to course",
		}
	} else {
		var courseUser models.CourseUser
		if err := initializers.DB.
			Where("course_id", courseId).
			Where("user_id", userData.ID).
			First(&courseUser).Error; err != nil {
			return &common.ErrorResponse{
				Code:    403,
				Message: "User does not belong to course",
				Details: err.Error(),
			}
		}

		for _, r := range courseUser.Roles {
			if r == role {
				return nil
			}
		}

		return &common.ErrorResponse{
			Code:    403,
			Message: "User does not belong to course",
		}
	}
}
