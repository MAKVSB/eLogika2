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
	switch userRole {
	case enums.CourseUserRoleAdmin:
		item, err := r.courseItemRepo.GetCourseItemByID(dbRef, courseID, courseItemID, userID, filters, full, version)
		if err != nil {
			item.Editable = true
		}
		return item, err
	case enums.CourseUserRoleGarant:
		item, err := r.courseItemRepo.GetCourseItemByID(dbRef, courseID, courseItemID, userID, filters, full, version)
		if err == nil {
			item.Editable = item.ManagedBy == enums.CourseUserRoleGarant
		}
		return item, err
	case enums.CourseUserRoleTutor:
		modifier := func(db *gorm.DB) *gorm.DB {
			if filters != nil {
				db = (*filters)(db)
			}
			return db.Where("managed_by = ? OR (managed_by = ? AND created_by_id = ?)", enums.CourseUserRoleGarant, enums.CourseUserRoleTutor, userID)
		}
		item, err := r.courseItemRepo.GetCourseItemByID(dbRef, courseID, courseItemID, userID, &modifier, full, version)
		if err == nil {
			item.Editable = item.ManagedBy == enums.CourseUserRoleTutor
		}
		return item, err
	default:
		return nil, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}
