package repositories

import (
	"encoding/json"
	"strconv"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/utils"
	"gorm.io/gorm"
)

type ChapterRepository struct{}

func NewChapterRepository() *ChapterRepository {
	return &ChapterRepository{}
}

func (r *ChapterRepository) GetChapterByID(
	dbRef *gorm.DB,
	courseId uint,
	chapterID uint,
	full bool,
	version *uint,
) (*models.Chapter, *common.ErrorResponse) {

	query := dbRef.
		Where("id = ?", chapterID).
		Where("course_id = ?", courseId)

	if full {
		query = query.Preload("Childs")
	}

	var chapter *models.Chapter
	if err := query.
		Find(&chapter).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    401,
			Message: "Failed to load chapter",
		}
	}

	if version != nil {
		if chapter.Version != *version {
			return nil, &common.ErrorResponse{
				Code:    409,
				Message: "Version mismatched",
				Details: strconv.Itoa(int(*version)) + " " + strconv.Itoa(int(chapter.Version)),
			}
		}
	}

	return chapter, nil
}

func (r *ChapterRepository) GetChapterByIDStudent(
	dbRef *gorm.DB,
	courseId uint,
	chapterID uint,
	full bool,
	version *uint,
) (*models.Chapter, *common.ErrorResponse) {

	query := dbRef.
		Where("id = ?", chapterID).
		Where("course_id = ?", courseId).
		Where("visible = ?", true)

	if full {
		query = query.Preload("Childs", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("visible = ?", true)
		})
	}

	var chapter *models.Chapter
	if err := query.
		Find(&chapter).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    401,
			Message: "Failed to load chapter",
		}
	}

	if version != nil {
		if chapter.Version != *version {
			return nil, &common.ErrorResponse{
				Code:    409,
				Message: "Version mismatched",
				Details: strconv.Itoa(int(*version)) + " " + strconv.Itoa(int(chapter.Version)),
			}
		}
	}

	return chapter, nil
}

func (r *ChapterRepository) ListChapters(
	dbRef *gorm.DB,
	courseID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Chapter, int64, *common.ErrorResponse) {
	query := dbRef.
		Model(models.Chapter{}).
		Where("course_id = ?", courseID)

	if full {
		query = query.Preload("Childs")
	}

	if filters != nil {
		query = (*filters)(query)
	}

	// Apply filters, sorting, pagination
	query, err := models.Chapter{}.ApplyFilters(query, searchParams.ColumnFilters, models.Chapter{}, map[string]interface{}{})
	if err != nil {
		return nil, 0, err
	}
	query = models.Chapter{}.ApplySorting(query, searchParams.Sorting)
	totalCount := models.Chapter{}.GetCount(query) // Gets count before pagination
	query = models.Chapter{}.ApplyPagination(query, searchParams.Pagination)

	var chapters []*models.Chapter
	if err := query.
		Find(&chapters).Error; err != nil {
		return nil, 0, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch chapter",
			Details: err.Error(),
		}
	}

	return chapters, totalCount, nil
}

func (r *ChapterRepository) ListChaptersStudent(
	dbRef *gorm.DB,
	courseID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Chapter, int64, *common.ErrorResponse) {
	query := dbRef.
		Model(models.Chapter{}).
		Where("course_id = ?", courseID).
		Where("visible = ?", true)

	if full {
		query = query.Preload("Childs", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("visible = ?", true)
		})
	}

	if filters != nil {
		query = (*filters)(query)
	}

	// Apply filters, sorting, pagination
	query, err := models.Chapter{}.ApplyFilters(query, searchParams.ColumnFilters, models.Chapter{}, map[string]interface{}{})
	if err != nil {
		return nil, 0, err
	}
	query = models.Chapter{}.ApplySorting(query, searchParams.Sorting)
	totalCount := models.Chapter{}.GetCount(query) // Gets count before pagination
	query = models.Chapter{}.ApplyPagination(query, searchParams.Pagination)

	var chapters []*models.Chapter
	if err := query.
		Find(&chapters).Error; err != nil {
		return nil, 0, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch chapter",
			Details: err.Error(),
		}
	}

	return chapters, totalCount, nil
}

func (r *ChapterRepository) SyncFiles(
	dbRef *gorm.DB,
	content json.RawMessage,
	chapter *models.Chapter,
) *common.ErrorResponse {
	if err := dbRef.Where("id IN ?", utils.GetFilesInsideContent(content)).Find(&chapter.ContentFiles).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to load files",
		}
	}

	if err := dbRef.Model(&chapter).Association("ContentFiles").Replace(&chapter.ContentFiles); err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update files",
		}
	}
	return nil
}

func (r *ChapterRepository) LastChapterOrder(
	dbRef *gorm.DB,
	parentId uint,
) uint {
	var lastChapter models.Chapter
	if err := dbRef.
		Where("parent_id = ?", parentId).
		Order("\"order\" DESC").
		First(&lastChapter).Error; err != nil {
		return 0
	}
	return lastChapter.Order
}

func (r *ChapterRepository) GetNextInOrder(
	dbRef *gorm.DB,
	courseId uint,
	parentId uint,
	direction enums.MoveDirectionEnum,
	originalOrder uint,
) (*models.Chapter, *common.ErrorResponse) {
	query := dbRef.
		Model(&models.Chapter{}).
		Where("course_id = ?", courseId).
		Where("parent_id = ?", parentId)

	if direction == enums.MoveDirectionDown {
		query = query.Where("\"order\" > ?", originalOrder).Order("\"order\" ASC")
	} else if direction == enums.MoveDirectionUp {
		query = query.Where("\"order\" < ?", originalOrder).Order("\"order\" DESC")
	} else {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Invalid move direction",
		}
	}

	var swapWithChapter *models.Chapter
	if err := query.First(&swapWithChapter).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch neighbor chapter",
		}
	}

	return swapWithChapter, nil
}
