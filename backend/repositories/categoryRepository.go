package repositories

import (
	"strconv"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"gorm.io/gorm"
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (r *CategoryRepository) GetCategoryByID(
	dbRef *gorm.DB,
	courseId uint,
	categoryID uint,
	version *uint,
) (*models.Category, *common.ErrorResponse) {

	var category *models.Category
	if err := dbRef.
		Preload("Steps").
		Preload("Chapter").
		Where("id = ?", categoryID).
		Where("course_id = ?", courseId).
		Find(&category).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    401,
			Message: "Failed to load category",
		}
	}

	if version != nil {
		if category.Version != *version {
			return nil, &common.ErrorResponse{
				Code:    409,
				Message: "Version mismatched",
				Details: strconv.Itoa(int(*version)) + " " + strconv.Itoa(int(category.Version)),
			}
		}
	}

	return category, nil
}

func (r *CategoryRepository) ListCategories(
	dbRef *gorm.DB,
	courseID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Category, int64, *common.ErrorResponse) {
	query := dbRef.
		Model(models.Category{}).
		Where("course_id = ?", courseID).
		Preload("Steps").
		Preload("Chapter")

	if filters != nil {
		query = (*filters)(query)
	}

	// Apply filters, sorting, pagination
	query, err := models.Question{}.ApplyFilters(query, searchParams.ColumnFilters, models.Question{}, map[string]interface{}{})
	if err != nil {
		return nil, 0, err
	}
	query = models.Question{}.ApplySorting(query, searchParams.Sorting)
	totalCount := models.Question{}.GetCount(query) // Gets count before pagination
	query = models.Question{}.ApplyPagination(query, searchParams.Pagination)

	var categories []*models.Category
	if err := query.
		Find(&categories).Error; err != nil {
		return nil, 0, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch category",
			Details: err.Error(),
		}
	}

	return categories, totalCount, nil
}
