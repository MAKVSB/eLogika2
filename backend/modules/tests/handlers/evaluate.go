package handlers

import (
	"fmt"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/tests/dtos"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestEvaluationRequest struct {
	InstanceID uint `uri:"instanceId" binding:"required"`
}

type EvaluateResponse struct {
	InstanceData dtos.TestInstanceDTO `json:"instanceData"`
}

// @Summary (Re)evaluates test
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param instanceId path int true "ID of the corresponding test instance"
// @Param body body TestEvaluationRequest true "Ability to filter results"
// @Success 200 {object} EvaluateResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/tests/evaluate [POST]

func TestEvaluate(c *gin.Context, userData authdtos.LoggedUserDTO) {
	// Load request data
	err, _, reqData := utils.GetRequestData[
		any,
		TestEvaluationRequest,
	](c)
	if err != nil {
		c.AbortWithStatusJSON(err.Code, err)
		return
	}

	// TODO validate from here

	// TODO check that user is in course, instance is active
	// // check permissions
	// coursePermissions := auth.GetClaimCourse(userData.Courses, params.CourseID)
	// if coursePermissions == nil || (!coursePermissions.IsTutor() &&
	// 	!coursePermissions.IsGarant() &&
	// 	!coursePermissions.IsAdmin()) {
	// 	c.JSON(403, common.ErrorResponse{
	// 		Message: "Not enough permissions",
	// 	})
	// 	return
	// }

	// if testInstance.ParticipantID != userData.ID {
	// 	c.AbortWithStatusJSON(401, common.ErrorResponse{
	// 		Message: "User is not allowed to access test instance",
	// 	})
	// 	return
	// }

	// if testInstance.IsExpired(time.Now()) {
	// 	c.AbortWithStatusJSON(401, common.ErrorResponse{
	// 		Code:    403,
	// 		Message: "Time expired",
	// 	})
	// 	return
	// }

	// if testInstance.State != enums.TestInstanceStateActive {
	// 	c.AbortWithStatusJSON(401, common.ErrorResponse{
	// 		Message: "Test already finished",
	// 	})
	// 	return
	// }

	// c.JSON(200, TestInstanceGetResponse{
	// 	InstanceData: dtos.TestInstanceDTO{}.From(testInstance),
	// })

	err = EvaluateTestInstance(initializers.DB, reqData.InstanceID, &userData)
	if err != nil {
		c.JSON(500, err)
		return
	}

	// services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
	// err = services_course_item.UpdateSelectedResults(transaction, params.CourseID, rootCoureItem, testInstance.ParticipantID)
	// if err != nil {
	// 	transaction.Rollback()
	// 	return err
	// }

	c.JSON(200, true)
}

func EvaluateTestInstance(dbRef *gorm.DB, instanceID uint, userData *authdtos.LoggedUserDTO) *common.ErrorResponse {
	var testInstance models.TestInstance
	if err := dbRef.
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("TestQuestion").
				Joins("TestQuestion.Question")
		}).
		Preload("Questions.Answers", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("TestQuestionAnswer").
				Joins("TestQuestionAnswer.Answer")
		}).
		Joins("Test").
		Joins("CourseItem").
		Joins("CourseItem.Parent").
		Joins("Result").
		Find(&testInstance, instanceID).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to load instance data",
			Details: err.Error(),
		}
	}
	if testInstance.State != enums.TestInstanceStateFinished && testInstance.State != enums.TestInstanceStateExpired {
		return &common.ErrorResponse{
			Code:    401,
			Message: "Cannot evaluate points for active test",
		}
	}
	type EvaluationQuestionAnswer struct {
		AllocatedPoints float64
		IsCorrect       bool
		IsSelected      bool
	}

	type EvaluationQuestion struct {
		AllocatedPoints      float64
		QuestionFormat       enums.QuestionFormatEnum
		Answers              []*EvaluationQuestionAnswer
		TextAnswerPercentage int
		TextAnswerReviewed   bool
	}

	type EvaluationBlock struct {
		AllocatedPoints       float64
		Weight                uint
		WrongAnswerPercentage uint
		Questions             []*EvaluationQuestion
	}
	// remap into a structure that is easier to process
	evaluationMap := make(map[uint]*EvaluationBlock)

	for _, block := range testInstance.Test.Blocks {
		evaluationMap[block.ID] = &EvaluationBlock{
			Weight:                block.Weight,
			WrongAnswerPercentage: block.WrongAnswerPercentage,
			Questions:             make([]*EvaluationQuestion, 0),
			AllocatedPoints:       float64(testInstance.CourseItem.PointsMax) * (float64(block.Weight) / 100),
		}
	}
	for _, q := range testInstance.Questions {
		blockId := q.TestQuestion.BlockID
		var evaluationAnswer EvaluationQuestion

		switch q.TestQuestion.Question.QuestionFormat {
		case enums.QuestionFormatOpen:
			evaluationAnswer = EvaluationQuestion{
				QuestionFormat:       q.TestQuestion.Question.QuestionFormat,
				TextAnswerPercentage: int(q.TextAnswerPercentage),
				TextAnswerReviewed:   q.TextAnswerReviewedByID != nil,
			}

		case enums.QuestionFormatTest:
			answers := make([]*EvaluationQuestionAnswer, len(q.Answers))

			for qa_i, qa := range q.Answers {
				answers[qa_i] = &EvaluationQuestionAnswer{
					IsCorrect:  qa.TestQuestionAnswer.Answer.Correct,
					IsSelected: qa.Selected,
				}
			}

			evaluationAnswer = EvaluationQuestion{
				QuestionFormat: q.TestQuestion.Question.QuestionFormat,
				Answers:        answers,
			}

		default:
			panic(fmt.Sprintf("unexpected enums.QuestionFormatEnum: %#v", q.TestQuestion.Question.QuestionFormat))
		}

		evaluationMap[blockId].Questions = append(evaluationMap[blockId].Questions, &evaluationAnswer)
	}
	//second loop over to assign allocated points correctly

	for _, ev_b := range evaluationMap {
		perQuestionPoints := ev_b.AllocatedPoints / float64(len(ev_b.Questions))
		for _, ev_b_q := range ev_b.Questions {
			ev_b_q.AllocatedPoints = perQuestionPoints
			if ev_b_q.QuestionFormat == enums.QuestionFormatTest {
				perAnswerPoints := ev_b_q.AllocatedPoints / float64(len(ev_b_q.Answers))
				for _, ev_b_q_a := range ev_b_q.Answers {
					ev_b_q_a.AllocatedPoints = perAnswerPoints
				}
			}
		}
	}
	// calculate points

	points := float64(0)
	final := true

	for _, block := range evaluationMap {
		for _, block_question := range block.Questions {
			switch block_question.QuestionFormat {
			case enums.QuestionFormatOpen:
				if !block_question.TextAnswerReviewed {
					final = false
				}
				points += (float64(block_question.TextAnswerPercentage) / 100) * block_question.AllocatedPoints
			case enums.QuestionFormatTest:
				for _, block_question_answer := range block_question.Answers {
					if block_question_answer.IsCorrect == block_question_answer.IsSelected {
						points += block_question_answer.AllocatedPoints
					} else {
						points -= block_question_answer.AllocatedPoints // * (float64(block.WrongAnswerPercentage) / 100)
					}
				}
			default:
				panic(fmt.Sprintf("unexpected enums.QuestionFormatEnum: %#v", block_question.QuestionFormat))
			}
		}
	}

	points += testInstance.BonusPoints

	testInstance.Result.Version = testInstance.Result.Version + 1
	testInstance.Result.Points = points
	testInstance.Result.Final = final
	if userData != nil {
		testInstance.Result.UpdatedByID = &userData.ID
	} else {
		testInstance.Result.UpdatedByID = nil
	}
	if err := dbRef.Save(&testInstance.Result).Error; err != nil {
		return &common.ErrorResponse{
			Message: "Failed to save result",
			Details: err.Error(),
		}
	}

	return nil
}
