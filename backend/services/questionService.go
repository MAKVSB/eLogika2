package services

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/questions/dtos"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/utils/tiptap"
	"gorm.io/gorm"
)

type QuestionService struct {
	questionRepo *repositories.QuestionRepository
}

func NewQuestionService(repo *repositories.QuestionRepository) *QuestionService {
	return &QuestionService{questionRepo: repo}
}

func (r *QuestionService) GetQuestionByID(
	dbRef *gorm.DB,
	courseID uint,
	questionID uint,
	userID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
) (*models.Question, *common.ErrorResponse) {
	switch userRole {
	case enums.CourseUserRoleAdmin:
		return r.questionRepo.GetQuestionByID(dbRef, courseID, questionID, userID, filters, full, version)
	case enums.CourseUserRoleGarant:
		modifier := func(db *gorm.DB) *gorm.DB {
			if filters != nil {
				db = (*filters)(db)
			}
			return db.Where("managed_by = ?", enums.CourseUserRoleGarant)
		}
		return r.questionRepo.GetQuestionByID(dbRef, courseID, questionID, userID, &modifier, full, version)
	case enums.CourseUserRoleTutor:
		modifier := func(db *gorm.DB) *gorm.DB {
			if filters != nil {
				db = (*filters)(db)
			}
			return db.Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userID)
		}
		return r.questionRepo.GetQuestionByID(dbRef, courseID, questionID, userID, &modifier, full, version)
	default:
		return nil, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}

func (r *QuestionService) ListQuestions(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Question, int64, *common.ErrorResponse) {
	switch userRole {
	case enums.CourseUserRoleAdmin:
		return r.questionRepo.ListQuestions(dbRef, courseID, userID, filters, full, searchParams)
	case enums.CourseUserRoleGarant:
		modifier := func(db *gorm.DB) *gorm.DB {
			if filters != nil {
				db = (*filters)(db)
			}
			return db.Where("managed_by = ?", enums.CourseUserRoleGarant)
		}
		return r.questionRepo.ListQuestions(dbRef, courseID, userID, &modifier, full, searchParams)
	case enums.CourseUserRoleTutor:
		modifier := func(db *gorm.DB) *gorm.DB {
			if filters != nil {
				db = (*filters)(db)
			}
			return db.Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userID)
		}
		return r.questionRepo.ListQuestions(dbRef, courseID, userID, &modifier, full, searchParams)
	default:
		return nil, 0, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}

func (r *QuestionService) SyncAnswers(
	dbRef *gorm.DB,
	userId uint,
	question *models.Question,
	answers []dtos.QuestionAnswerAdminDTO,
) *common.ErrorResponse {
	answerIds := make([]uint, 0)
	for _, answer := range answers {
		var answerData models.Answer
		if answer.ID == 0 {
			answerData = models.Answer{
				ID:      0,
				Version: 1,
			}
		} else {
			if err := dbRef.First(&answerData, answer.ID).Error; err != nil {
				return &common.ErrorResponse{
					Code:    404,
					Message: "Failed to load answer",
				}
			}
			answerData.Version = answer.Version + 1
		}

		err := tiptap.FindAndSaveRelations(dbRef, userId, answer.Content, &answerData, "ContentFiles")
		if err != nil {
			return err
		}
		answerData.Content = answer.Content

		err = tiptap.FindAndSaveRelations(dbRef, userId, answer.Explanation, &answerData, "ExplanationFiles")
		if err != nil {
			return err
		}
		answerData.Explanation = answer.Explanation

		answerData.TimeToSolve = answer.TimeToSolve
		answerData.Correct = answer.Correct

		if err := dbRef.Save(&answerData).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to create or update answer",
				Details: answer,
			}
		}

		// If new answer, update link to question
		if answer.ID == 0 {
			question_answer := models.QuestionAnswer{
				Version:    1,
				QuestionID: question.ID,
				AnswerID:   answerData.ID,
			}

			if err := dbRef.Save(&question_answer).Error; err != nil {
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to link answer to question",
					Details: answer,
				}
			}
		}

		answerIds = append(answerIds, answerData.ID)
	}
	{
		var itemsToDelete []models.QuestionAnswer
		dq := dbRef.
			Model(&models.QuestionAnswer{}).
			Where("question_id = ?", question.ID)

		if len(answerIds) != 0 {
			dq = dq.Where("answer_id NOT IN (?)", answerIds)
		}

		if err := dq.Scan(&itemsToDelete).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to fetch answers to delete data",
			}
		}

		for _, item := range itemsToDelete {
			questionUsed, err := r.questionRepo.IsAnswerUsedByTest(dbRef, item.ID)
			if err != nil {
				return err
			}
			if questionUsed {
				return &common.ErrorResponse{
					Code:    409,
					Message: "Any of deleted answers have been already used by test.",
					Details: "If you want to delete answer for future tests, create a new version",
				}
			}

			if err := dbRef.
				Delete(&item).Error; err != nil {
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to unlink answer from question",
				}
			}
		}
	}
	return nil
}
