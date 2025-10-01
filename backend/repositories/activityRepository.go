package repositories

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"gorm.io/gorm"
)

type ActivityRepository struct{}

func NewActivityRepository() *ActivityRepository {
	return &ActivityRepository{}
}

// func (r *ActivityRepository) GetActivityByID(
// 	transaction *gorm.DB,
// 	courseID uint,
// 	activityID uint,
// 	userID uint,
// 	filters *(func(*gorm.DB) *gorm.DB),
// 	full bool,
// ) (*models.Activity, *common.ErrorResponse) {
// 	query := transaction

// 	if filters != nil {
// 		query = (*filters)(query)
// 	}

// 	if full {
// 		query = query.
// 			Preload("CourseItem").
// 			Preload("Term").
// 			Preload("CreatedBy").
// 			Preload("Blocks").
// 			Preload("Questions").
// 			Preload("Instances")
// 	}

// 	var activity *models.Activity
// 	if err := query.
// 		First(&activity, activityID).Error; err != nil {
// 		transaction.Rollback()
// 		return nil, &common.ErrorResponse{
// 			Code:    404,
// 			Message: "Failed to fetch activity",
// 			Details: err.Error(),
// 		}
// 	}

// 	return activity, nil
// }

// // // Modifications for ease of writing code later
// // func (r *ActivityRepository) GetActivityByIDAdmin(
// // 	transaction *gorm.DB,
// // 	courseID uint,
// // 	activityID uint,
// // 	userID uint,
// // 	full bool,
// // 	version *uint,
// // ) (*models.Activity, *common.ErrorResponse) {
// // 	return r.GetActivityByID(transaction, courseID, activityID, userID, nil, full, version)
// // }

// // func (r *ActivityRepository) GetActivityByIDGarant(
// // 	transaction *gorm.DB,
// // 	courseID uint,
// // 	activityID uint,
// // 	userID uint,
// // 	full bool,
// // 	version *uint,
// // ) (*models.Activity, *common.ErrorResponse) {
// // 	modifier := func(db *gorm.DB) *gorm.DB {
// // 		return db.Where("managed_by = ?", enums.CourseUserRoleGarant)
// // 	}
// // 	return r.GetActivityByID(transaction, courseID, activityID, userID, &modifier, full, version)
// // }

// // func (r *ActivityRepository) GetActivityByIDTutor(
// // 	transaction *gorm.DB,
// // 	courseID uint,
// // 	activityID uint,
// // 	userID uint,
// // 	full bool,
// // 	version *uint,
// // ) (*models.Activity, *common.ErrorResponse) {
// // 	modifier := func(db *gorm.DB) *gorm.DB {
// // 		return db.Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userID)
// // 	}
// // 	return r.GetActivityByID(transaction, courseID, activityID, userID, &modifier, full, version)
// // }

func (r *ActivityRepository) GetActivityInstanceByID(
	dbRef *gorm.DB,
	activityInstanceID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	courseItemID *uint,
	participantID *uint,
) (*models.ActivityInstance, *common.ErrorResponse) {
	if courseItemID == nil && participantID == nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Instance ownership must be check using CourseItem or Participant",
		}
	}

	query := dbRef

	if courseItemID != nil {
		query = query.Where("course_item_id = ?", *courseItemID)
	}

	if participantID != nil {
		query = query.Where("participant_id = ?", *participantID)
	}

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
		query = query.
			InnerJoins("CourseItem").
			InnerJoins("CourseItem.ActivityDetail").
			InnerJoins("Result")
	}

	var activityInstance *models.ActivityInstance
	if err := query.
		First(&activityInstance, activityInstanceID).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch activity",
			Details: err.Error(),
		}
	}

	return activityInstance, nil
}

func (r *ActivityRepository) ListActivityInstances(
	dbRef *gorm.DB,
	courseItemID uint,
	termID *uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.ActivityInstance, int64, *common.ErrorResponse) {
	query := dbRef.
		Model(&models.ActivityInstance{}).
		Preload("Participant").
		Where("course_item_id = ?", courseItemID)

	if termID != nil && *termID != 0 {
		query = query.Where("term_id = ?", *termID)
	}

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
	}

	// Apply filters, sorting, pagination
	query, err := models.ActivityInstance{}.ApplyFilters(query, searchParams.ColumnFilters, models.ActivityInstance{}, map[string]interface{}{}, "")
	if err != nil {
		return nil, 0, err
	}
	query = models.ActivityInstance{}.ApplySorting(query, searchParams.Sorting)
	totalCount := models.ActivityInstance{}.GetCount(query) // Gets count before pagination
	query = models.ActivityInstance{}.ApplyPagination(query, searchParams.Pagination)

	var activitys []*models.ActivityInstance
	if err := query.
		Find(&activitys).Error; err != nil {
		return nil, 0, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch activitys",
			Details: err.Error(),
		}
	}

	return activitys, totalCount, nil
}
