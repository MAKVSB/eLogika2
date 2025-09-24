package services

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"gorm.io/gorm"
)

type ActivityService struct {
	activityRepo *repositories.ActivityRepository
}

func NewActivityService(repo *repositories.ActivityRepository) *ActivityService {
	return &ActivityService{activityRepo: repo}
}

// func (r *ActivityService) GetActivityInstanceByID(
// 	dbRef *gorm.DB,
// 	courseID uint,
// 	activityID uint,
// 	userID uint,
// 	userRole enums.CourseUserRoleEnum,
// 	filters *(func(*gorm.DB) *gorm.DB),
// 	full bool,
// 	version *uint,
// ) (*models.Activity, *common.ErrorResponse) {
// 	if userRole == enums.CourseUserRoleAdmin {
// 		return r.activityRepo.GetActivityByIDAdmin(dbRef, courseID, activityID, userID, full, version)
// 	} else if userRole == enums.CourseUserRoleGarant {
// 		return r.activityRepo.GetActivityByIDGarant(dbRef, courseID, activityID, userID, full, version)
// 	} else if userRole == enums.CourseUserRoleTutor {
// 		return r.activityRepo.GetActivityByIDTutor(dbRef, courseID, activityID, userID, full, version)
// 	} else {
// 		return nil, &common.ErrorResponse{
// 			Code:    403,
// 			Message: "Not enough permissions",
// 		}
// 	}
// }

func (r *ActivityService) ListActivityInstances(
	dbRef *gorm.DB,
	courseItemID uint,
	termID *uint,
	userID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.ActivityInstance, int64, *common.ErrorResponse) {
	if userRole == enums.CourseUserRoleAdmin {
		return r.activityRepo.ListActivityInstances(dbRef, courseItemID, termID, userID, filters, full, searchParams)
	} else if userRole == enums.CourseUserRoleGarant {
		return r.activityRepo.ListActivityInstances(dbRef, courseItemID, termID, userID, filters, full, searchParams)
	} else if userRole == enums.CourseUserRoleTutor {
		return r.activityRepo.ListActivityInstances(dbRef, courseItemID, termID, userID, filters, full, searchParams)
	} else {
		return nil, 0, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}
