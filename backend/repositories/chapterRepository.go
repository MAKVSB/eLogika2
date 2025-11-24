package repositories

import (
	"strconv"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
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
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
	onlyVisible bool,
) (*models.Chapter, *common.ErrorResponse) {

	query := dbRef.
		Where("id = ?", chapterID).
		Where("course_id = ?", courseId)

	if onlyVisible {
		query = query.Where("visible = ?", true)
	}

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
		query = query.Preload("Childs", func(db *gorm.DB) *gorm.DB {
			if onlyVisible {
				db = db.Where("visible = ?", true)
			}
			return db
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
	onlyVisible bool,
) ([]*models.Chapter, int64, *common.ErrorResponse) {
	var err *common.ErrorResponse
	query := dbRef.
		Model(models.Chapter{}).
		Where("course_id = ?", courseID)

	if onlyVisible {
		query = query.Where("visible = ?", true)
	}

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
		query = query.Preload("Childs", func(db *gorm.DB) *gorm.DB {
			if onlyVisible {
				db = db.Where("visible = ?", true)
			}
			return db
		})
	}

	// Apply filters, sorting, pagination
	if searchParams != nil {
		query, err = models.Chapter{}.ApplyFilters(query, searchParams.ColumnFilters, models.Chapter{}, map[string]interface{}{}, "")
		if err != nil {
			return nil, 0, err
		}
		query = models.Chapter{}.ApplySorting(query, searchParams.Sorting, "id DESC")
	}
	totalCount := models.Chapter{}.GetCount(query) // Gets count before pagination
	if searchParams != nil {
		query = models.Chapter{}.ApplyPagination(query, searchParams.Pagination)
	}

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

	switch direction {
	case enums.MoveDirectionDown:
		query = query.Where("\"order\" > ?", originalOrder).Order("\"order\" ASC")
	case enums.MoveDirectionUp:
		query = query.Where("\"order\" < ?", originalOrder).Order("\"order\" DESC")
	default:
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
