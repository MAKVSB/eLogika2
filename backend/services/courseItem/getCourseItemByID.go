package services_course_item

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

func (r *CourseItemService) GetCourseItemByID(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	userID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
) (*models.CourseItem, *common.ErrorResponse) {
	if userRole == enums.CourseUserRoleAdmin {
		item, err := r.courseItemRepo.GetCourseItemByIDAdmin(dbRef, courseID, courseItemID, userID, full, version)
		if err != nil {
			item.Editable = true
		}
		return item, err
	} else if userRole == enums.CourseUserRoleGarant {
		item, err := r.courseItemRepo.GetCourseItemByIDGarant(dbRef, courseID, courseItemID, userID, full, version)
		if err == nil {
			item.Editable = item.ManagedBy == enums.CourseUserRoleGarant
		}
		return item, err
	} else if userRole == enums.CourseUserRoleTutor {
		item, err := r.courseItemRepo.GetCourseItemByIDTutor(dbRef, courseID, courseItemID, userID, full, version)
		if err == nil {
			item.Editable = item.ManagedBy == enums.CourseUserRoleTutor
		}
		return item, err
	} else {
		return nil, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}
