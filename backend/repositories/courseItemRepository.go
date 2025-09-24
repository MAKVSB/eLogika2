package repositories

import (
	"strconv"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type CourseItemRepository struct{}

func NewCourseItemRepository() *CourseItemRepository {
	return &CourseItemRepository{}
}

func (r *CourseItemRepository) GetCourseItemByID(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
) (*models.CourseItem, *common.ErrorResponse) {
	query := dbRef.
		Where("course_id = ?", courseID)

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
		query = query.Preload("ActivityDetail").
			Preload("TestDetail").
			Preload("GroupDetail").
			Preload("Children").
			Preload("Terms").
			Preload("Parent")
	}

	var courseItem *models.CourseItem
	if err := query.
		First(&courseItem, courseItemID).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch courseItem",
			Details: err.Error(),
		}
	}

	if version != nil {
		if courseItem.Version != *version {
			return nil, &common.ErrorResponse{
				Code:    409,
				Message: "Version mismatched",
				Details: strconv.Itoa(int(*version)) + " " + strconv.Itoa(int(courseItem.Version)),
			}
		}
	}

	return courseItem, nil
}

// Modifications for ease of writing code later
func (r *CourseItemRepository) GetCourseItemByIDAdmin(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	userID uint,
	full bool,
	version *uint,
) (*models.CourseItem, *common.ErrorResponse) {
	return r.GetCourseItemByID(dbRef, courseID, courseItemID, userID, nil, full, version)
}

func (r *CourseItemRepository) GetCourseItemByIDGarant(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	userID uint,
	full bool,
	version *uint,
) (*models.CourseItem, *common.ErrorResponse) {
	return r.GetCourseItemByID(dbRef, courseID, courseItemID, userID, nil, full, version)
}

func (r *CourseItemRepository) GetCourseItemByIDTutor(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	userID uint,
	full bool,
	version *uint,
) (*models.CourseItem, *common.ErrorResponse) {
	modifier := func(db *gorm.DB) *gorm.DB {
		return db.Where("managed_by = ? OR (managed_by = ? AND created_by_id = ?)", enums.CourseUserRoleGarant, enums.CourseUserRoleTutor, userID)
	}
	return r.GetCourseItemByID(dbRef, courseID, courseItemID, userID, &modifier, full, version)
}

func (r *CourseItemRepository) ListCourseItems(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.CourseItem, int64, *common.ErrorResponse) {
	query := dbRef.
		Model(&models.CourseItem{}).
		Where("course_id = ?", courseID)

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
		query = query.Preload("ActivityDetail").
			Preload("TestDetail").
			Preload("GroupDetail").
			Preload("Children").
			Preload("Terms").
			Preload("Parent").
			Preload("Results")
	}

	// Apply filters, sorting, pagination
	query, err := models.CourseItem{}.ApplyFilters(query, searchParams.ColumnFilters, models.CourseItem{}, map[string]interface{}{})
	if err != nil {
		return nil, 0, err
	}
	query = models.CourseItem{}.ApplySorting(query, searchParams.Sorting)
	totalCount := models.CourseItem{}.GetCount(query) // Gets count before pagination
	query = models.CourseItem{}.ApplyPagination(query, searchParams.Pagination)

	var courseItems []*models.CourseItem
	if err := query.
		Find(&courseItems).Error; err != nil {
		return nil, 0, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch courseItem",
			Details: err.Error(),
		}
	}

	return courseItems, totalCount, nil
}

func (r *CourseItemRepository) ListCourseItemsAdmin(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.CourseItem, int64, *common.ErrorResponse) {
	return r.ListCourseItems(dbRef, courseID, userID, nil, full, searchParams)
}

func (r *CourseItemRepository) ListCourseItemsGarant(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.CourseItem, int64, *common.ErrorResponse) {
	modifier := func(db *gorm.DB) *gorm.DB {
		return db.Where("managed_by = ?", enums.CourseUserRoleGarant)
	}
	return r.ListCourseItems(dbRef, courseID, userID, &modifier, full, searchParams)
}

func (r *CourseItemRepository) ListCourseItemsTutor(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.CourseItem, int64, *common.ErrorResponse) {
	modifier := func(db *gorm.DB) *gorm.DB {
		return db.Where("managed_by = ? OR (managed_by = ? AND created_by_id = ?)", enums.CourseUserRoleGarant, enums.CourseUserRoleTutor, userID)
	}
	return r.ListCourseItems(dbRef, courseID, userID, &modifier, full, searchParams)
}
