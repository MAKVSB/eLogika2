package repositories

import (
	"strconv"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type TemplateRepository struct{}

func NewTemplateRepository() *TemplateRepository {
	return &TemplateRepository{}
}

func (r *TemplateRepository) GetTemplateByID(
	dbRef *gorm.DB,
	courseID uint,
	templateID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
) (*models.Template, *common.ErrorResponse) {
	query := dbRef.
		Where("course_id = ?", courseID).
		InnerJoins("CreatedBy")

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
		query = query.
			Preload("Blocks").
			Preload("Blocks.Segments").
			Preload("Blocks.Segments.Questions").
			Preload("Blocks.Segments.Steps")
	}

	var template *models.Template
	if err := query.
		First(&template, templateID).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch template",
			Details: err.Error(),
		}
	}

	if version != nil {
		if template.Version != *version {
			return nil, &common.ErrorResponse{
				Code:    409,
				Message: "Version mismatched",
				Details: strconv.Itoa(int(*version)) + " " + strconv.Itoa(int(template.Version)),
			}
		}
	}

	return template, nil
}

// Modifications for ease of writing code later
func (r *TemplateRepository) GetTemplateByIDAdmin(
	dbRef *gorm.DB,
	courseID uint,
	templateID uint,
	userID uint,
	full bool,
	version *uint,
) (*models.Template, *common.ErrorResponse) {
	return r.GetTemplateByID(dbRef, courseID, templateID, userID, nil, full, version)
}

func (r *TemplateRepository) GetTemplateByIDGarant(
	dbRef *gorm.DB,
	courseID uint,
	templateID uint,
	userID uint,
	full bool,
	version *uint,
) (*models.Template, *common.ErrorResponse) {
	modifier := func(db *gorm.DB) *gorm.DB {
		return db.Where("managed_by = ?", enums.CourseUserRoleGarant)
	}
	return r.GetTemplateByID(dbRef, courseID, templateID, userID, &modifier, full, version)
}

func (r *TemplateRepository) GetTemplateByIDTutor(
	dbRef *gorm.DB,
	courseID uint,
	templateID uint,
	userID uint,
	full bool,
	version *uint,
) (*models.Template, *common.ErrorResponse) {
	modifier := func(db *gorm.DB) *gorm.DB {
		return db.Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userID)
	}
	return r.GetTemplateByID(dbRef, courseID, templateID, userID, &modifier, full, version)
}

func (r *TemplateRepository) ListTemplates(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Template, int64, *common.ErrorResponse) {
	query := dbRef.
		Model(&models.Template{}).
		Where("course_id = ?", courseID).
		InnerJoins("CreatedBy").
		Preload("Blocks")

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
		query = query.
			Preload("Blocks.Segments").
			Preload("Blocks.Segments.Questions").
			Preload("Blocks.Segments.Steps")
	}

	// Apply filters, sorting, pagination
	query, err := models.Template{}.ApplyFilters(query, searchParams.ColumnFilters, models.Template{}, map[string]interface{}{})
	if err != nil {
		return nil, 0, err
	}
	query = models.Template{}.ApplySorting(query, searchParams.Sorting)
	totalCount := models.Template{}.GetCount(query) // Gets count before pagination
	query = models.Template{}.ApplyPagination(query, searchParams.Pagination)

	var templates []*models.Template
	if err := query.
		Find(&templates).Error; err != nil {
		return nil, 0, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch template",
			Details: err.Error(),
		}
	}

	return templates, totalCount, nil
}

func (r *TemplateRepository) ListTemplatesAdmin(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Template, int64, *common.ErrorResponse) {
	return r.ListTemplates(dbRef, courseID, userID, nil, full, searchParams)
}

func (r *TemplateRepository) ListTemplatesGarant(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Template, int64, *common.ErrorResponse) {
	modifier := func(db *gorm.DB) *gorm.DB {
		return db.Where("managed_by = ?", enums.CourseUserRoleGarant)
	}
	return r.ListTemplates(dbRef, courseID, userID, &modifier, full, searchParams)
}

func (r *TemplateRepository) ListTemplatesTutor(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Template, int64, *common.ErrorResponse) {
	modifier := func(db *gorm.DB) *gorm.DB {
		return db.Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userID)
	}
	return r.ListTemplates(dbRef, courseID, userID, &modifier, full, searchParams)
}
