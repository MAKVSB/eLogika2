package services

import (
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"gorm.io/gorm"
)

type ClassService struct {
	classRepo *repositories.ClassRepository
}

func NewClassService(repo *repositories.ClassRepository) *ClassService {
	return &ClassService{classRepo: repo}
}

func (r *ClassService) GetClassByID(
	dbRef *gorm.DB,
	courseID uint,
	classID uint,
	userID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
) (*models.Class, *common.ErrorResponse) {
	switch userRole {
	case enums.CourseUserRoleAdmin:
		return r.classRepo.GetClassByID(dbRef, courseID, classID, userID, filters, full, version)
	case enums.CourseUserRoleGarant:
		return r.classRepo.GetClassByID(dbRef, courseID, classID, userID, filters, full, version)
	case enums.CourseUserRoleTutor:
		modifier := func(db *gorm.DB) *gorm.DB {
			if filters != nil {
				db = (*filters)(db)
			}
			return db.Joins("inner join class_tutors as help1 on help1.class_id = classes.id AND help1.user_id = ? AND help1.deleted_at is NULL", userID)
		}
		return r.classRepo.GetClassByID(dbRef, courseID, classID, userID, &modifier, full, version)
	default:
		return nil, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}

func (r *ClassService) ListClasses(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Class, int64, *common.ErrorResponse) {
	switch userRole {
	case enums.CourseUserRoleAdmin:
		return r.classRepo.ListClasses(initializers.DB, courseID, userID, filters, true, searchParams)
	case enums.CourseUserRoleGarant:
		return r.classRepo.ListClasses(initializers.DB, courseID, userID, filters, true, searchParams)
	case enums.CourseUserRoleTutor:
		modifier := func(db *gorm.DB) *gorm.DB {
			if filters != nil {
				db = (*filters)(db)
			}
			return db.Joins("inner join class_tutors as help1 on help1.class_id = classes.id AND help1.user_id = ? AND help1.deleted_at is NULL", userID)
		}
		return r.classRepo.ListClasses(initializers.DB, courseID, userID, &modifier, true, searchParams)
	default:
		return nil, 0, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

}
