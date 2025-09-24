package services_course_item

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

func (r *CourseItemService) ListCourseItems(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.CourseItem, int64, *common.ErrorResponse) {
	if userRole == enums.CourseUserRoleAdmin {
		items, itemcount, err := r.courseItemRepo.ListCourseItemsAdmin(dbRef, courseID, userID, full, searchParams)
		for _, item := range items {
			item.Editable = true
		}
		return items, itemcount, err
	} else if userRole == enums.CourseUserRoleGarant {
		items, itemcount, err := r.courseItemRepo.ListCourseItemsGarant(dbRef, courseID, userID, full, searchParams)
		for _, item := range items {
			item.Editable = item.ManagedBy == enums.CourseUserRoleGarant
		}
		return items, itemcount, err
	} else if userRole == enums.CourseUserRoleTutor {
		items, itemcount, err := r.courseItemRepo.ListCourseItemsTutor(dbRef, courseID, userID, full, searchParams)
		for _, item := range items {
			item.Editable = item.ManagedBy == enums.CourseUserRoleTutor
		}
		return items, itemcount, err
	} else {
		return nil, 0, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}
