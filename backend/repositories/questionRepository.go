package repositories

import (
	"strconv"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type QuestionRepository struct{}

func NewQuestionRepository() *QuestionRepository {
	return &QuestionRepository{}
}

func (r *QuestionRepository) GetQuestionByID(
	dbRef *gorm.DB,
	courseID uint,
	questionID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
) (*models.Question, *common.ErrorResponse) {
	query := dbRef.
		InnerJoins("CourseLink", initializers.DB.Where("CourseLink.course_id = ?", courseID))

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
		query = query.
			Preload("CourseLink.Steps").
			InnerJoins("CreatedBy").
			Preload("CheckedBy").
			Preload("CheckedBy.User").
			Preload("Answers").
			Preload("Answers.Answer")
	}

	var question *models.Question
	if err := query.
		First(&question, questionID).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch question",
			Details: err.Error(),
		}
	}

	if version != nil {
		if question.Version != *version {
			return nil, &common.ErrorResponse{
				Code:    409,
				Message: "Version mismatched",
				Details: strconv.Itoa(int(*version)) + " " + strconv.Itoa(int(question.Version)),
			}
		}
	}

	return question, nil
}

// Modifications for ease of writing code later
func (r *QuestionRepository) GetQuestionByIDAdmin(
	dbRef *gorm.DB,
	courseID uint,
	questionID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
) (*models.Question, *common.ErrorResponse) {
	return r.GetQuestionByID(dbRef, courseID, questionID, userID, filters, full, version)
}

func (r *QuestionRepository) GetQuestionByIDGarant(
	dbRef *gorm.DB,
	courseID uint,
	questionID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
) (*models.Question, *common.ErrorResponse) {
	modifier := func(db *gorm.DB) *gorm.DB {
		if filters != nil {
			db = (*filters)(db)
		}
		return db.Where("managed_by = ?", enums.CourseUserRoleGarant)
	}
	return r.GetQuestionByID(dbRef, courseID, questionID, userID, &modifier, full, version)
}

func (r *QuestionRepository) GetQuestionByIDTutor(
	dbRef *gorm.DB,
	courseID uint,
	questionID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
) (*models.Question, *common.ErrorResponse) {
	modifier := func(db *gorm.DB) *gorm.DB {
		if filters != nil {
			db = (*filters)(db)
		}
		return db.Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userID)
	}
	return r.GetQuestionByID(dbRef, courseID, questionID, userID, &modifier, full, version)
}

func (r *QuestionRepository) GetMaxVersion(
	dbRef *gorm.DB,
	questionGroupID uint,
) (uint, *common.ErrorResponse) {
	var maxVersion uint
	if err := dbRef.
		Raw("SELECT MAX(version) FROM questions WHERE question_group_id = ?", questionGroupID).
		Scan(&maxVersion).Error; err != nil {
		return 0, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to get max version",
			Details: err.Error(),
		}
	}
	return maxVersion, nil
}

func (r *QuestionRepository) SyncSteps(
	dbRef *gorm.DB,
	question *models.Question,
	categoryId *uint,
	newSteps []uint,
) (*models.Question, *common.ErrorResponse) {
	if err := dbRef.
		Model(&models.Step{}).
		Where("id IN ? AND category_id = ?", newSteps, categoryId).
		Find(&question.CourseLink.Steps).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to load steps",
		}
	}

	if err := dbRef.
		Model(&question.CourseLink).
		Association("Steps").
		Replace(&question.CourseLink.Steps); err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update steps",
		}
	}

	return question, nil
}

func (r *QuestionRepository) IsAnswerUsedByTest(
	dbRef *gorm.DB,
	answerId uint,
) (bool, *common.ErrorResponse) {
	var isUsed uint
	if err := dbRef.
		Raw("SELECT 1 FROM test_question_answers WHERE answer_id = ?", answerId).
		Scan(&isUsed).Error; err != nil {
		return isUsed != 0, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to get max version",
			Details: err.Error(),
		}
	}
	return isUsed != 0, nil
}

func (r *QuestionRepository) ListQuestions(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Question, int64, *common.ErrorResponse) {
	courseLinkQuery := initializers.DB.Where("CourseLink.course_id = ?", courseID)

	// Apply filters to innerobjects
	courseLinkQuery, err := models.CourseQuestion{}.ApplyFilters(courseLinkQuery, searchParams.ColumnFilters, models.CourseQuestion{}, map[string]interface{}{}, "CourseLink.")
	if err != nil {
		return nil, 0, err
	}

	query := dbRef.
		Model(&models.Question{}).
		InnerJoins("CourseLink", courseLinkQuery).
		InnerJoins("CreatedBy").
		Preload("CheckedBy").
		Preload("CheckedBy.User")

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
		query = query.
			Preload("CourseLink.Steps").
			Preload("Answers").
			Preload("Answers.Answer")
	}

	// Apply filters, sorting, pagination
	query, err = models.Question{}.ApplyFilters(query, searchParams.ColumnFilters, models.Question{}, map[string]interface{}{
		"userID": userID,
	})
	if err != nil {
		return nil, 0, err
	}
	query = models.Question{}.ApplySorting(query, searchParams.Sorting)
	totalCount := models.Question{}.GetCount(query) // Gets count before pagination
	query = models.Question{}.ApplyPagination(query, searchParams.Pagination)

	var questions []*models.Question
	if err := query.
		Find(&questions).Error; err != nil {
		return nil, 0, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch question",
			Details: err.Error(),
		}
	}

	return questions, totalCount, nil
}

func (r *QuestionRepository) ListQuestionsAdmin(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Question, int64, *common.ErrorResponse) {
	return r.ListQuestions(dbRef, courseID, userID, filters, full, searchParams)
}

func (r *QuestionRepository) ListQuestionsGarant(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Question, int64, *common.ErrorResponse) {
	modifier := func(db *gorm.DB) *gorm.DB {
		if filters != nil {
			db = (*filters)(db)
		}
		return db.Where("managed_by = ?", enums.CourseUserRoleGarant)
	}
	return r.ListQuestions(dbRef, courseID, userID, &modifier, full, searchParams)
}

func (r *QuestionRepository) ListQuestionsTutor(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Question, int64, *common.ErrorResponse) {
	modifier := func(db *gorm.DB) *gorm.DB {
		if filters != nil {
			db = (*filters)(db)
		}
		return db.Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userID)
	}
	return r.ListQuestions(dbRef, courseID, userID, &modifier, full, searchParams)
}
