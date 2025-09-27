package services

import (
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
		return r.classRepo.GetClassByIDAdmin(dbRef, courseID, classID, userID, filters, full, version)
	case enums.CourseUserRoleGarant:
		return r.classRepo.GetClassByIDGarant(dbRef, courseID, classID, userID, filters, full, version)
	case enums.CourseUserRoleTutor:
		return r.classRepo.GetClassByIDTutor(dbRef, courseID, classID, userID, filters, full, version)
	default:
		return nil, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}
