package auth

import (
	"elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
)

func GetClaimCourse(courses []dtos.LoggedUserCourseDTO, courseId uint) *dtos.LoggedUserCourseDTO {
	for _, item := range courses {
		if item.ID == courseId {
			return &item
		}
	}
	return nil
}

func GetClaimCourseRole(courses []dtos.LoggedUserCourseDTO, courseId uint, role enums.CourseUserRoleEnum) *common.ErrorResponse {
	for _, item := range courses {
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
}
