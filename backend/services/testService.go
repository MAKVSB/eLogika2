package services

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"gorm.io/gorm"
)

type TestService struct {
	testRepo *repositories.TestRepository
}

func NewTestService(repo *repositories.TestRepository) *TestService {
	return &TestService{testRepo: repo}
}

// func (r *TestService) GetTestByID(
// 	dbRef *gorm.DB,
// 	courseID uint,
// 	testID uint,
// 	userID uint,
// 	userRole enums.CourseUserRoleEnum,
// 	filters *(func(*gorm.DB) *gorm.DB),
// 	full bool,
// 	version *uint,
// ) (*models.Test, *common.ErrorResponse) {
// 	if userRole == enums.CourseUserRoleAdmin {
// 		return r.testRepo.GetTestByIDAdmin(dbRef, courseID, testID, userID, full, version)
// 	} else if userRole == enums.CourseUserRoleGarant {
// 		return r.testRepo.GetTestByIDGarant(dbRef, courseID, testID, userID, full, version)
// 	} else if userRole == enums.CourseUserRoleTutor {
// 		return r.testRepo.GetTestByIDTutor(dbRef, courseID, testID, userID, full, version)
// 	} else {
// 		return nil, &common.ErrorResponse{
// 			Code:    403,
// 			Message: "Not enough permissions",
// 		}
// 	}
// }

func (r *TestService) ListTests(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	termID *uint,
	userID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Test, int64, *common.ErrorResponse) {
	switch userRole {
	case enums.CourseUserRoleAdmin:
		return r.testRepo.ListTests(dbRef, courseID, courseItemID, termID, userID, filters, full, searchParams)
	case enums.CourseUserRoleGarant:
		return r.testRepo.ListTests(dbRef, courseID, courseItemID, termID, userID, filters, full, searchParams)
	case enums.CourseUserRoleTutor:
		return r.testRepo.ListTests(dbRef, courseID, courseItemID, termID, userID, filters, full, searchParams)
	default:
		return nil, 0, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}

func (r *TestService) ListTestInstances(
	dbRef *gorm.DB,
	courseItemID uint,
	termID *uint,
	testID uint,
	userID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.TestInstance, int64, *common.ErrorResponse) {
	switch userRole {
	case enums.CourseUserRoleAdmin:
		return r.testRepo.ListTestInstances(dbRef, courseItemID, termID, testID, userID, filters, full, searchParams)
	case enums.CourseUserRoleGarant:
		return r.testRepo.ListTestInstances(dbRef, courseItemID, termID, testID, userID, filters, full, searchParams)
	case enums.CourseUserRoleTutor:
		return r.testRepo.ListTestInstances(dbRef, courseItemID, termID, testID, userID, filters, full, searchParams)
	default:
		return nil, 0, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}
