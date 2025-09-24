package repositories

import (
	"strconv"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/utils"
	"gorm.io/gorm"
)

type TermRepository struct{}

func NewTermRepository() *TermRepository {
	return &TermRepository{}
}

func (r *TermRepository) GetTermByID(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	termID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
) (*models.Term, *common.ErrorResponse) {
	query := dbRef.
		Where("course_id = ?", courseID).
		Where("course_item_id = ?", courseItemID)

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
	}

	var term *models.Term
	if err := query.
		First(&term, termID).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch term",
			Details: err.Error(),
		}
	}

	if version != nil {
		if term.Version != *version {
			return nil, &common.ErrorResponse{
				Code:    409,
				Message: "Version mismatched",
				Details: strconv.Itoa(int(*version)) + " " + strconv.Itoa(int(term.Version)),
			}
		}
	}

	return term, nil
}

// Modifications for ease of writing code later
func (r *TermRepository) GetTermByIDAdmin(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	termID uint,
	userID uint,
	full bool,
	version *uint,
) (*models.Term, *common.ErrorResponse) {
	return r.GetTermByID(dbRef, courseID, courseItemID, termID, userID, nil, full, version)
}

func (r *TermRepository) GetTermByIDGarant(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	termID uint,
	userID uint,
	full bool,
	version *uint,
) (*models.Term, *common.ErrorResponse) {
	// TODO zkontrolovat, že může editovat courseitem
	return r.GetTermByID(dbRef, courseID, courseItemID, termID, userID, nil, full, version)
}

func (r *TermRepository) GetTermByIDTutor(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	termID uint,
	userID uint,
	full bool,
	version *uint,
) (*models.Term, *common.ErrorResponse) {
	// TODO zkontrolovat, že může editovat courseitem
	// modifier := func(db *gorm.DB) *gorm.DB {
	// 	return db.Where("managed_by = ? OR (managed_by = ? AND created_by_id = ?)", enums.CourseUserRoleGarant, enums.CourseUserRoleTutor, userID)
	// }
	return r.GetTermByID(dbRef, courseID, courseItemID, termID, userID, nil, full, version)
}

func (r *TermRepository) ListTerms(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Term, int64, *common.ErrorResponse) {
	query := dbRef.
		Model(&models.Term{}).
		Where("course_id = ?", courseID).
		Where("course_item_id = ?", courseItemID).
		Preload("Students")

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
	}

	// Apply filters, sorting, pagination
	query, err := models.Term{}.ApplyFilters(query, searchParams.ColumnFilters, models.Term{}, map[string]interface{}{})
	if err != nil {
		return nil, 0, err
	}
	query = models.Term{}.ApplySorting(query, searchParams.Sorting)
	totalCount := models.Term{}.GetCount(query) // Gets count before pagination
	query = models.Term{}.ApplyPagination(query, searchParams.Pagination)

	var terms []*models.Term
	if err := query.
		Find(&terms).Error; err != nil {
		return nil, 0, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch term",
			Details: err.Error(),
		}
	}

	return terms, totalCount, nil
}

func (r *TermRepository) ListTermsAdmin(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	userID uint,
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Term, int64, *common.ErrorResponse) {
	return r.ListTerms(dbRef, courseID, courseItemID, userID, nil, full, searchParams)
}

func (r *TermRepository) ListTermsGarant(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	userID uint,
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Term, int64, *common.ErrorResponse) {
	// TODO zkontrolovat, že může editovat courseitem
	return r.ListTerms(dbRef, courseID, courseItemID, userID, nil, full, searchParams)
}

func (r *TermRepository) ListTermsTutor(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	userID uint,
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Term, int64, *common.ErrorResponse) {
	// TODO zkontrolovat, že může editovat courseitem
	// modifier := func(db *gorm.DB) *gorm.DB {
	// 	return db.Where("managed_by = ? OR (managed_by = ? AND created_by_id = ?)", enums.CourseUserRoleGarant, enums.CourseUserRoleTutor, userID)
	// }
	return r.ListTerms(dbRef, courseID, courseItemID, userID, nil, full, searchParams)
}

func (r *TermRepository) ListJoinedStudents(
	dbRef *gorm.DB,
	termID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.UserTerm, int64, *common.ErrorResponse) {
	query := dbRef.
		Model(&models.UserTerm{}).
		Preload("User").
		Where("term_id = ?", termID)

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
	}

	utils.DebugPrintJSON(searchParams)

	// Apply filters, sorting, pagination
	if searchParams != nil {
		var err *common.ErrorResponse
		query, err = models.UserTerm{}.ApplyFilters(query, searchParams.ColumnFilters, models.UserTerm{}, map[string]interface{}{})
		if err != nil {
			return nil, 0, err
		}
		query = models.UserTerm{}.ApplySorting(query, searchParams.Sorting)
	}
	totalCount := models.UserTerm{}.GetCount(query) // Gets count before pagination
	if searchParams != nil {
		query = models.UserTerm{}.ApplyPagination(query, searchParams.Pagination)
	}

	var terms []*models.UserTerm
	if err := query.
		Find(&terms).Error; err != nil {
		return nil, 0, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch joined students",
			Details: err.Error(),
		}
	}

	return terms, totalCount, nil
}
