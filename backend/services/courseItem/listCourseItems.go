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
	switch userRole {
	case enums.CourseUserRoleAdmin:
		items, itemcount, err := r.courseItemRepo.ListCourseItems(dbRef, courseID, userID, filters, full, searchParams)
		for _, item := range items {
			item.Editable = true
		}
		return items, itemcount, err
	case enums.CourseUserRoleGarant:
		modifier := func(db *gorm.DB) *gorm.DB {
			if filters != nil {
				db = (*filters)(db)
			}
			return db.Where("managed_by = ?", enums.CourseUserRoleGarant)
		}
		items, itemcount, err := r.courseItemRepo.ListCourseItems(dbRef, courseID, userID, &modifier, full, searchParams)

		for _, item := range items {
			item.Editable = item.ManagedBy == enums.CourseUserRoleGarant
		}
		return items, itemcount, err
	case enums.CourseUserRoleTutor:
		modifier := func(db *gorm.DB) *gorm.DB {
			if filters != nil {
				db = (*filters)(db)
			}
			return db.Where("managed_by = ? OR (managed_by = ? AND created_by_id = ?)", enums.CourseUserRoleGarant, enums.CourseUserRoleTutor, userID)
		}
		items, itemcount, err := r.courseItemRepo.ListCourseItems(dbRef, courseID, userID, &modifier, full, searchParams)

		for _, item := range items {
			item.Editable = item.ManagedBy == enums.CourseUserRoleTutor
		}
		return items, itemcount, err
	default:
		return nil, 0, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}
