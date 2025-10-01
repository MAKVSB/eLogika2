package repositories

import (
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type TestRepository struct{}

func NewTestRepository() *TestRepository {
	return &TestRepository{}
}

func (r *TestRepository) GetTestByID(
	transaction *gorm.DB,
	courseID uint,
	testID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
) (*models.Test, *common.ErrorResponse) {
	query := transaction

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
		query = query.
			Preload("CourseItem").
			Preload("Term").
			Preload("CreatedBy").
			Preload("Blocks").
			Preload("Questions").
			Preload("Instances")
	}

	var test *models.Test
	if err := query.
		First(&test, testID).Error; err != nil {
		transaction.Rollback()
		return nil, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch test",
			Details: err.Error(),
		}
	}

	return test, nil
}

// // Modifications for ease of writing code later
// func (r *TestRepository) GetTestByIDAdmin(
// 	transaction *gorm.DB,
// 	courseID uint,
// 	testID uint,
// 	userID uint,
// 	full bool,
// 	version *uint,
// ) (*models.Test, *common.ErrorResponse) {
// 	return r.GetTestByID(transaction, courseID, testID, userID, nil, full, version)
// }

// func (r *TestRepository) GetTestByIDGarant(
// 	transaction *gorm.DB,
// 	courseID uint,
// 	testID uint,
// 	userID uint,
// 	full bool,
// 	version *uint,
// ) (*models.Test, *common.ErrorResponse) {
// 	modifier := func(db *gorm.DB) *gorm.DB {
// 		return db.Where("managed_by = ?", enums.CourseUserRoleGarant)
// 	}
// 	return r.GetTestByID(transaction, courseID, testID, userID, &modifier, full, version)
// }

// func (r *TestRepository) GetTestByIDTutor(
// 	transaction *gorm.DB,
// 	courseID uint,
// 	testID uint,
// 	userID uint,
// 	full bool,
// 	version *uint,
// ) (*models.Test, *common.ErrorResponse) {
// 	modifier := func(db *gorm.DB) *gorm.DB {
// 		return db.Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userID)
// 	}
// 	return r.GetTestByID(transaction, courseID, testID, userID, &modifier, full, version)
// }

func (r *TestRepository) GetTestInstanceByID(
	dbRef *gorm.DB,
	testInstanceID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	isTutor bool,
	courseItemID *uint,
	participantID *uint,
) (*models.TestInstance, *common.ErrorResponse) {
	if courseItemID == nil && participantID == nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Instance ownership must be check using CourseItem or Participant",
		}
	}

	query := dbRef

	if courseItemID != nil {
		query = query.Where("test_instances.course_item_id = ?", *courseItemID)
	}

	if participantID != nil {
		query = query.Where("participant_id = ?", *participantID)
	}

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
		query = query.
			InnerJoins("Participant").
			InnerJoins("Term").
			InnerJoins("Result").
			InnerJoins("CourseItem").
			InnerJoins("CourseItem.TestDetail").
			Preload("Questions", func(db *gorm.DB) *gorm.DB {
				return db.
					Unscoped().
					Joins("TestQuestion", initializers.DB.Unscoped()).
					Joins("TestQuestion.Question", initializers.DB.Unscoped()).
					Order("TestQuestion__order ASC")
			}).
			Preload("Questions.Answers", func(db *gorm.DB) *gorm.DB {
				return db.
					Unscoped().
					Joins("TestQuestionAnswer", initializers.DB.Unscoped()).
					Joins("TestQuestionAnswer.Answer", initializers.DB.Unscoped()).
					Order("TestQuestionAnswer__order ASC")
			})
	}

	var testInstance *models.TestInstance
	if err := query.
		First(&testInstance, testInstanceID).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch test",
			Details: err.Error(),
		}
	}
	if isTutor {
		for i, question := range testInstance.Questions {
			if question.TestQuestion.Question.QuestionFormat == enums.QuestionFormatOpen {
				var openAnswers []models.Answer
				if err := initializers.DB.
					InnerJoins("Question", initializers.DB.Where("question_id = ?", question.TestQuestion.QuestionID)).
					Find(&openAnswers).Error; err != nil {
					return nil, &common.ErrorResponse{
						Code:    500,
						Message: "Failed to load open question data",
					}
				}
				testInstance.Questions[i].TestQuestion.OpenAnswers = openAnswers
			}
		}
	}

	return testInstance, nil
}

func (r *TestRepository) ListTests(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	termID *uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Test, int64, *common.ErrorResponse) {
	query := dbRef.
		Model(&models.Test{}).
		InnerJoins("CreatedBy").
		Where("course_id = ?", courseID).
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
	query, err := models.Test{}.ApplyFilters(query, searchParams.ColumnFilters, models.Test{}, map[string]interface{}{}, "")
	if err != nil {
		return nil, 0, err
	}
	query = models.Test{}.ApplySorting(query, searchParams.Sorting)
	totalCount := models.Test{}.GetCount(query) // Gets count before pagination
	query = models.Test{}.ApplyPagination(query, searchParams.Pagination)

	var tests []*models.Test
	if err := query.
		Find(&tests).Error; err != nil {
		return nil, 0, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch tests",
			Details: err.Error(),
		}
	}

	return tests, totalCount, nil
}

func (r *TestRepository) ListTestsAdmin(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	termID *uint,
	userID uint,
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Test, int64, *common.ErrorResponse) {
	return r.ListTests(dbRef, courseID, courseItemID, termID, userID, nil, full, searchParams)
}

func (r *TestRepository) ListTestsGarant(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	termID *uint,
	userID uint,
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Test, int64, *common.ErrorResponse) {
	return r.ListTests(dbRef, courseID, courseItemID, termID, userID, nil, full, searchParams)
}

func (r *TestRepository) ListTestsTutor(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	termID *uint,
	userID uint,
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Test, int64, *common.ErrorResponse) {
	return r.ListTests(dbRef, courseID, courseItemID, termID, userID, nil, full, searchParams)
}

func (r *TestRepository) ListTestInstances(
	dbRef *gorm.DB,
	courseItemID uint,
	termID *uint,
	testID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.TestInstance, int64, *common.ErrorResponse) {
	query := dbRef.
		Model(&models.TestInstance{}).
		Preload("Participant").
		Where("course_item_id = ?", courseItemID).
		Where("test_id = ?", testID)

	if termID != nil && *termID != 0 {
		query = query.Where("term_id = ?", *termID)
	}

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
	}

	// Apply filters, sorting, pagination
	query, err := models.TestInstance{}.ApplyFilters(query, searchParams.ColumnFilters, models.TestInstance{}, map[string]interface{}{}, "")
	if err != nil {
		return nil, 0, err
	}
	query = models.TestInstance{}.ApplySorting(query, searchParams.Sorting)
	totalCount := models.TestInstance{}.GetCount(query) // Gets count before pagination
	query = models.TestInstance{}.ApplyPagination(query, searchParams.Pagination)

	var tests []*models.TestInstance
	if err := query.
		Find(&tests).Error; err != nil {
		return nil, 0, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch tests",
			Details: err.Error(),
		}
	}

	return tests, totalCount, nil
}

func (r *TestRepository) ListTestInstancesAdmin(
	dbRef *gorm.DB,
	courseItemID uint,
	termID *uint,
	testID uint,
	userID uint,
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.TestInstance, int64, *common.ErrorResponse) {
	return r.ListTestInstances(dbRef, courseItemID, termID, testID, userID, nil, full, searchParams)
}

func (r *TestRepository) ListTestInstancesGarant(
	dbRef *gorm.DB,
	courseItemID uint,
	termID *uint,
	testID uint,
	userID uint,
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.TestInstance, int64, *common.ErrorResponse) {
	return r.ListTestInstances(dbRef, courseItemID, termID, testID, userID, nil, full, searchParams)
}

func (r *TestRepository) ListTestInstancesTutor(
	dbRef *gorm.DB,
	courseItemID uint,
	termID *uint,
	testID uint,
	userID uint,
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.TestInstance, int64, *common.ErrorResponse) {
	return r.ListTestInstances(dbRef, courseItemID, termID, testID, userID, nil, full, searchParams)
}
