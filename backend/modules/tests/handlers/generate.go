package handlers

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"strconv"
	"time"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/tests/helpers"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestGeneratorRequest struct {
	Variants *uint                      `json:"variants"`
	UsersAll bool                       `json:"usersAll"`
	UsersIDs []uint                     `json:"usersIds"`
	Form     enums.TestInstanceFormEnum `json:"form"  binding:"required"`
}

type TestGeneratorResponse struct {
	Success bool `json:"success"`
}

// @Summary Generate tests based on input data
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TestGeneratorRequest true "Ability to filter results"
// @Success 200 {object} TestGeneratorResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/tests/generate [post]
func Generate(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
			TermID       uint `uri:"termId" binding:"required"`
		},
		TestGeneratorRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		return err
	}
	var courseItem *models.CourseItem
	// Check if tutor/garant can view/modify courseItem
	if userRole == enums.CourseUserRoleAdmin {
		if err := initializers.DB.
			Preload("TestDetail").
			Find(&courseItem, params.CourseItemID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permission for this item",
			}
		}
	} else if userRole == enums.CourseUserRoleGarant {
		if err := initializers.DB.
			Preload("TestDetail").
			Where("managed_by = ?", enums.CourseUserRoleGarant).
			Find(&courseItem, params.CourseItemID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permission for this item",
			}
		}
	} else if userRole == enums.CourseUserRoleTutor {
		if err := initializers.DB.
			Preload("TestDetail").
			Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userData.ID).
			Find(&courseItem, params.CourseItemID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permission for this item",
			}
		}
	} else {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	var template *models.Template
	query := initializers.DB.
		Preload("Blocks").
		Preload("Blocks.Segments").
		Preload("Blocks.Segments.Questions").
		Preload("Blocks.Segments.Steps").
		Where("course_id = ?", params.CourseID)

	if err := query.First(&template, courseItem.TestDetail.TestTemplateID).Error; err != nil {
		return &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch template",
			Details: err.Error(),
		}
	}

	generatorCache, err := helpers.LoadQuestionsByTemplate(template, courseItem)
	if err != nil {
		return err
	}

	transaction := initializers.DB.Begin()

	if reqData.Variants != nil && *reqData.Variants != 0 {
		generatedTestVariants := make([]*models.Test, 0)

		for var_i := range int(*reqData.Variants) {
			variant, err := GenerateTest(
				transaction,
				template,
				generatorCache,
				params.CourseID,
				params.CourseItemID,
				params.TermID,
				&userData,
				userData.Username,
				GetVariantLabel(var_i),
			)
			if err != nil {
				transaction.Rollback()
				return err
			}
			generatedTestVariants = append(generatedTestVariants, variant)
		}

		// Save tests to database
		batchSize := 50
		total := len(generatedTestVariants)

		for i := 0; i < total; i += batchSize {
			end := i + batchSize
			if end > total {
				end = total
			}

			batch := generatedTestVariants[i:end]

			if err := transaction.Save(&batch).Error; err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to save generated test",
					Details: err.Error(),
				}
			}
		}
	} else if reqData.UsersAll || len(reqData.UsersIDs) != 0 {
		// Load all joined students
		var joinedUsers []*models.UserTerm
		query := transaction.
			Preload("User").
			Where("term_id = ?", params.TermID)

		if len(reqData.UsersIDs) != 0 {
			query = query.Where("user_id in ?", reqData.UsersIDs)
		}

		if err := query.Find(&joinedUsers).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    404,
				Message: "Failed to fetch joined users",
				Details: err.Error(),
			}
		}

		generatedTestInstances := make([]*models.TestInstance, 0)

		for _, ju := range joinedUsers {
			generatedTest, err := GenerateTest(
				transaction,
				template,
				generatorCache,
				params.CourseID,
				params.CourseItemID,
				params.TermID,
				&userData,
				ju.User.Username,
				"",
			)
			if err != nil {
				transaction.Rollback()
				return err
			}

			if err := transaction.Save(&generatedTest).Error; err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to save generated test",
					Details: err.Error(),
				}
			}

			testInstance, err := helpers.CreateInstance(
				transaction,
				generatedTest,
				ju.UserID,
				params.TermID,
				courseItem.ID,
				reqData.Form,
			)
			if err != nil {
				transaction.Rollback()
				return err
			}

			generatedTestInstances = append(generatedTestInstances, testInstance)
		}

		// Save tests to database
		batchSize := 50
		total := len(generatedTestInstances)

		for i := 0; i < total; i += batchSize {
			end := i + batchSize
			if end > total {
				end = total
			}

			batch := generatedTestInstances[i:end]

			if err := transaction.Save(&batch).Error; err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to save generated instances",
					Details: err.Error(),
				}
			}
		}
	} else {
		return &common.ErrorResponse{
			Code:    404,
			Message: "One of available inputs must be specified",
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	c.JSON(200, TestGeneratorResponse{
		Success: true,
	})

	return nil
}

func GenerateSingleTest(
	dbRef *gorm.DB,
	courseId uint,
	templateId uint,
	termId uint,
	generatingByUser *authdtos.LoggedUserDTO,
	generatingForUser string,
	courseItemId uint,
) (*models.Test, *common.ErrorResponse) {
	var template *models.Template
	query := initializers.DB.
		Preload("Blocks").
		Preload("Blocks.Segments").
		Preload("Blocks.Segments.Questions").
		Preload("Blocks.Segments.Steps").
		Where("course_id = ?", courseId)

	if err := query.First(&template, templateId).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch template",
			Details: err.Error(),
		}
	}

	var courseItem *models.CourseItem
	if err := dbRef.First(&courseItem, courseItemId).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch course item",
			Details: err.Error(),
		}
	}

	generatorCache, err := helpers.LoadQuestionsByTemplate(template, courseItem)
	if err != nil {
		return nil, err
	}

	generatedTestVariant, err := GenerateTest(
		dbRef,
		template,
		generatorCache,
		courseId,
		courseItemId,
		termId,
		generatingByUser,
		generatingForUser,
		"",
	)
	if err != nil {
		return nil, err
	}

	// Save tests to database
	if err := dbRef.Save(&generatedTestVariant).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to save generated test",
			Details: err.Error(),
		}
	}

	return generatedTestVariant, nil
}

func GenerateTest(
	transaction *gorm.DB,
	template *models.Template,
	generatorCache *helpers.GeneratorCache,
	courseId uint,
	courseItemId uint,
	termId uint,
	generatingByUser *authdtos.LoggedUserDTO,
	generatingForUser string,
	group string,
) (*models.Test, *common.ErrorResponse) {
	var generatedTestVariant *models.Test

	maxTries := 3
	for try := range maxTries {
		testQuestions, err := GenerateVariantQuestions(generatorCache, template.MixBlocks, template.MixEverything)
		if err != nil {
			if try == maxTries-1 {
				return nil, &common.ErrorResponse{
					Code:    500,
					Message: "Failed to generate test variants",
					Details: err.Error(),
				}
			}
			continue
		}

		variant := models.Test{
			CourseID:     courseId,
			Name:         generatingForUser + "-" + strconv.FormatInt(time.Now().UnixMicro(), 10),
			CreatedByID:  generatingByUser.ID,
			Blocks:       make([]models.TestBlock, len(template.Blocks)),
			Questions:    testQuestions,
			CourseItemID: courseItemId,
			TermID:       termId,
			Group:        group,
		}

		for tb_i, tb := range template.Blocks {
			variant.Blocks[tb_i] = models.TestBlock{
				ID:                    tb.ID,
				Title:                 tb.Title,
				ShowName:              tb.ShowName,
				Weight:                tb.Weight,
				WrongAnswerPercentage: tb.WrongAnswerPercentage,
			}
		}

		generatedTestVariant = &variant
		break
	}

	return generatedTestVariant, nil
}

func GenerateVariantQuestions(generatorCache *helpers.GeneratorCache, mixBlocks bool, mixEverything bool) ([]models.TestQuestion, error) {
	variantQuestionIDs := []uint{}

	blockedQuestions := make([][]models.TestQuestion, 0)

	order := uint(1)

	for _, block := range generatorCache.Blocks {
		blockQuestions := make([]models.TestQuestion, 0)

		for s_id, segment := range block.Segments {
			Shuffle(segment.QuestionPool)
			SortQuestionCandidates(segment.QuestionPool)

			succesfullyPickedQuestionCount := 0
			for q := 0; q < len(segment.QuestionPool); q++ {
				if succesfullyPickedQuestionCount < int(segment.ReqQuestionCount) {
					if slices.ContainsFunc(variantQuestionIDs, func(qID uint) bool {
						return qID == segment.QuestionPool[q].ID
					}) {
						continue
					}

					pickedQuestion := models.TestQuestion{
						Order:      order,
						BlockID:    block.BlockData.ID,
						QuestionID: segment.QuestionPool[q].ID,
					}

					if segment.QuestionPool[q].QuestionFormat == enums.QuestionFormatTest {
						pickedAnswers, err := PickRandomAnswers(int(block.BlockData.AnswerCount), segment.QuestionPool[q].Answers, block.BlockData.AnswerDistribution)
						if err != nil {
							continue
						}
						pickedQuestion.Answers = *pickedAnswers
					}

					order++

					succesfullyPickedQuestionCount++
					variantQuestionIDs = append(variantQuestionIDs, segment.QuestionPool[q].ID)
					blockQuestions = append(blockQuestions, pickedQuestion)
				} else {
					break
				}
			}

			if succesfullyPickedQuestionCount != int(segment.ReqQuestionCount) {
				return nil, errors.New("not enough questions in question pool for segment " + strconv.Itoa(s_id))
			}
		}

		if block.BlockData.MixInsideBlock {
			Shuffle(blockQuestions)
		}

		blockedQuestions = append(blockedQuestions, blockQuestions)
	}

	if mixBlocks {
		Shuffle(blockedQuestions)
	}

	variantQuestions := make([]models.TestQuestion, 0)

	for _, questionsInBlock := range blockedQuestions {
		variantQuestions = append(variantQuestions, questionsInBlock...)
	}

	if mixEverything {
		Shuffle(variantQuestions)
	}

	return variantQuestions, nil
}

func PickRandomAnswers(reqAnswerCount int, allAnswers []helpers.QuestionAnswer, distribution enums.AnswerDistributionEnum) (*[]models.TestQuestionAnswer, error) {
	pickedAnswers := make([]models.TestQuestionAnswer, 0)

	// Prepare searching data
	var remainingIncorrect int
	var remainingCorrect int
	var remainingRandom int
	order := uint(1)

	switch distribution {
	case enums.AnswerDistributionExactlyOneCorrect:
		remainingIncorrect = max(reqAnswerCount-1, 0)
		remainingCorrect = 1
		remainingRandom = 0
	case enums.AnswerDistributionMinimumOneCorrect:
		remainingIncorrect = 0
		remainingCorrect = 1
		remainingRandom = max(reqAnswerCount-1, 0)
	case enums.AnswerDistributionMinimumOneCorrectOneIncorrect:
		remainingIncorrect = 1
		remainingCorrect = 1
		remainingRandom = max(reqAnswerCount-2, 0)
	default:
		panic(fmt.Sprintf("unexpected enums.AnswerDistributionEnum: %#v", distribution))
	}

	//Shuffle all answers
	Shuffle(allAnswers)

	//Loop over and pick
	for _, a := range allAnswers {
		if a.Answer.Correct {
			if remainingCorrect > 0 {
				remainingCorrect -= 1
				pickedAnswers = append(pickedAnswers, models.TestQuestionAnswer{
					AnswerID: a.Answer.ID,
					Order:    order,
				})
				order++
			} else if remainingRandom > 0 {
				remainingRandom -= 1
				pickedAnswers = append(pickedAnswers, models.TestQuestionAnswer{
					AnswerID: a.Answer.ID,
					Order:    order,
				})
				order++
			}
		} else {
			if remainingIncorrect > 0 {
				remainingIncorrect -= 1
				pickedAnswers = append(pickedAnswers, models.TestQuestionAnswer{
					AnswerID: a.Answer.ID,
					Order:    order,
				})
				order++
			} else if remainingRandom > 0 {
				remainingRandom -= 1
				pickedAnswers = append(pickedAnswers, models.TestQuestionAnswer{
					AnswerID: a.Answer.ID,
					Order:    order,
				})
				order++
			}
		}

		if remainingCorrect == 0 && remainingIncorrect == 0 && remainingRandom == 0 {
			break
		}
	}

	if remainingCorrect == 0 && remainingIncorrect == 0 && remainingRandom == 0 {

		Shuffle(pickedAnswers)

		return &pickedAnswers, nil
	} else {
		fmt.Println(remainingCorrect, remainingIncorrect, remainingRandom)
		return nil, errors.New("not enough answers to pick from")
	}
}

func Shuffle[T any](arr []T) []T {
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	return arr
}

func SortQuestionCandidates(arr []*helpers.TMPQ) []*helpers.TMPQ {
	slices.SortFunc(arr, func(a, b *helpers.TMPQ) int {
		return int(a.TimesUsed - b.TimesUsed)
	})
	return arr
}

func GetVariantLabel(n int) string {
	label := ""
	for n >= 0 {
		label = string(rune('A'+(n%26))) + label
		n = n/26 - 1
	}
	return label
}
