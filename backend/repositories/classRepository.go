package repositories

import (
	"strconv"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"gorm.io/gorm"
)

type ClassRepository struct{}

func NewClassRepository() *ClassRepository {
	return &ClassRepository{}
}

func (r *ClassRepository) GetClassByID(
	dbRef *gorm.DB,
	courseID uint,
	classID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
) (*models.Class, *common.ErrorResponse) {
	query := dbRef.
		Where("id = ?", classID).
		Where("course_id = ?", courseID).
		Preload("Tutors", func(db *gorm.DB) *gorm.DB {
			return db.
				InnerJoins("User")
		})

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
		query = query.
			Preload("Students", func(db *gorm.DB) *gorm.DB {
				return db.
					InnerJoins("User")
			}).
			Preload("Tutors", func(db *gorm.DB) *gorm.DB {
				return db.
					InnerJoins("User")
			})
	}

	var class *models.Class
	if err := query.
		First(&class, classID).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch class",
			Details: err.Error(),
		}
	}

	if version != nil {
		if class.Version != *version {
			return nil, &common.ErrorResponse{
				Code:    409,
				Message: "Version mismatched",
				Details: strconv.Itoa(int(*version)) + " " + strconv.Itoa(int(class.Version)),
			}
		}
	}

	return class, nil
}

func (r *ClassRepository) ListClasses(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Class, int64, *common.ErrorResponse) {
	var err *common.ErrorResponse
	query := dbRef.
		Model(&models.Class{}).
		Where("course_id = ?", courseID).
		Preload("Tutors", func(db *gorm.DB) *gorm.DB {
			return db.
				InnerJoins("User")
		})

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
		query = query.
			Preload("Students", func(db *gorm.DB) *gorm.DB {
				return db.
					InnerJoins("User")
			}).
			Preload("Tutors", func(db *gorm.DB) *gorm.DB {
				return db.
					InnerJoins("User")
			})
	}

	// Apply filters, sorting, pagination
	if searchParams != nil {
		query, err = models.Class{}.ApplyFilters(query, searchParams.ColumnFilters, models.Class{}, map[string]interface{}{}, "")
		if err != nil {
			return nil, 0, err
		}
		query = models.Class{}.ApplySorting(query, searchParams.Sorting, "id DESC")
	}
	totalCount := models.Class{}.GetCount(query) // Gets count before pagination
	if searchParams != nil {
		query = models.Class{}.ApplyPagination(query, searchParams.Pagination)
	}

	var classs []*models.Class
	if err := query.
		Find(&classs).Error; err != nil {
		return nil, 0, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch class",
			Details: err.Error(),
		}
	}

	return classs, totalCount, nil
}

func (r *ClassRepository) ListClassStudents(
	dbRef *gorm.DB,
	courseID uint,
	classID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.User, int64, *common.ErrorResponse) {
	var err *common.ErrorResponse
	query := dbRef.
		Model(&models.User{}).
		InnerJoins("INNER JOIN class_students ON class_students.user_id = users.id AND class_id = ?", classID)

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
	}

	// Apply filters, sorting, pagination
	if searchParams != nil {
		query, err = models.User{}.ApplyFilters(query, searchParams.ColumnFilters, models.User{}, map[string]interface{}{}, "")
		if err != nil {
			return nil, 0, err
		}
		query = models.User{}.ApplySorting(query, searchParams.Sorting, "")
	}
	totalCount := models.User{}.GetCount(query) // Gets count before pagination
	if searchParams != nil {
		query = models.User{}.ApplyPagination(query, searchParams.Pagination)
	}

	var classStudents []*models.User
	if err := query.
		Find(&classStudents).Error; err != nil {
		return nil, 0, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch class students",
			Details: err.Error(),
		}
	}

	return classStudents, totalCount, nil
}
