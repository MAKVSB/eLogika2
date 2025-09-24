package helpers

import (
	"fmt"
	"slices"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type Answer struct {
	ID          uint `gorm:"primarykey"`
	TimeToSolve int  `` // Time it takes to check the correctness of the answer
	Correct     bool `` // If answer is true
}

type QuestionAnswer struct {
	ID uint `gorm:"primarykey"`

	QuestionID uint ``
	AnswerID   uint ``

	Answer *Answer ``
}

type TMPQ struct {
	ID             uint
	TimesUsed      uint
	Answers        []QuestionAnswer `gorm:"foreignKey:QuestionID"`
	QuestionFormat enums.QuestionFormatEnum
}

type GeneratorCacheBlockSegment struct {
	ReqQuestionCount uint
	QuestionPool     []*TMPQ
}

type GeneratorCacheBlock struct {
	BlockData *models.TemplateBlock
	Segments  []GeneratorCacheBlockSegment
}

type GeneratorCache struct {
	Blocks []GeneratorCacheBlock
}

func LoadQuestionsByTemplate(template *models.Template, courseItem *models.CourseItem) (*GeneratorCache, *common.ErrorResponse) {
	generatorCache := GeneratorCache{}

	globalQuestionQuery := initializers.DB.Model(models.Question{}).Select("questions.*")
	globalQuestionQuery = globalQuestionQuery.Where("active = ?", true)
	globalQuestionQuery = globalQuestionQuery.InnerJoins("CourseLink", initializers.DB.Where("CourseLink.course_id = ?", template.CourseID))
	if courseItem.ManagedBy == enums.CourseUserRoleGarant {
		// Course item owned by garant, so picking only questions managed by garant
		globalQuestionQuery = globalQuestionQuery.Where("managed_by = ?", enums.CourseUserRoleGarant)
	} else {
		// Owned by tutor so picking only his questions
		globalQuestionQuery = globalQuestionQuery.Where("managed_by = ?", enums.CourseUserRoleTutor)
		globalQuestionQuery = globalQuestionQuery.Where("created_by_id = ?", courseItem.CreatedById)
	}

	// Load all possible questions into memory
	for _, blockData := range template.Blocks {
		blockQuestionQuery := globalQuestionQuery.Session(&gorm.Session{})
		// Select only questions from the same course as template
		// Question format
		blockQuestionQuery = blockQuestionQuery.Where("question_format = ?", blockData.QuestionFormat)
		// Only questions where is enough answers available
		// TODO Check question ownership (tutor/garant mode)

		// If not an open formatted question, load available answers
		if blockData.QuestionFormat == enums.QuestionFormatTest {
			blockQuestionQuery = blockQuestionQuery.Preload("Answers", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "question_id", "answer_id")
			}).Preload("Answers.Answer", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "time_to_solve", "correct")
			})
			blockQuestionQuery = blockQuestionQuery.Where("answer_count > ?", blockData.AnswerCount)
		}
		// TODO Requested question difficulty
		// blockQuestionQuery = blockQuestionQuery.Where("difficulty BETWEEN ? AND ?", blockData.difficultyFrom, blockData.diffitultyTo).

		generatorCacheBlock := GeneratorCacheBlock{
			BlockData: &blockData,
			Segments:  make([]GeneratorCacheBlockSegment, 0),
		}
		for _, segmentData := range blockData.Segments {
			segmentQuery := blockQuestionQuery.Session(&gorm.Session{})

			// Filter by chapter
			if segmentData.ChapterID != nil {
				segmentQuery = segmentQuery.Where("CourseLink.chapter_id = ?", *segmentData.ChapterID)

				// Filter by category
				if segmentData.CategoryID != nil {
					segmentQuery = segmentQuery.Where("CourseLink.category_id = ?", segmentData.CategoryID)
				}
			}

			var pickedQuestions []*TMPQ
			// Filter by switch
			switch segmentData.FilterBy {
			case enums.CategoryFilterALL:
				if err := segmentQuery.Model(models.Question{}).Find(&pickedQuestions).Error; err != nil {
					return nil, &common.ErrorResponse{
						Code:    404,
						Message: "Failed to fetch questions",
						Details: err.Error(),
					}
				}
			case enums.CategoryFilterQ:
				// Add hand picked questions
				handpickedQuery := blockQuestionQuery.Session(&gorm.Session{})
				segmentQuestionGroups := make([]uint, len(segmentData.Questions))
				for q_ind, q := range segmentData.Questions {
					segmentQuestionGroups[q_ind] = q.ID
				}
				if err := handpickedQuery.Where("question_group_id IN (?)", segmentQuestionGroups).Find(&pickedQuestions).Error; err != nil {
					return nil, &common.ErrorResponse{
						Code:    404,
						Message: "Failed to fetch questions",
						Details: err.Error(),
					}
				}
			case enums.CategoryFilterS:
				var allQuestions []models.Question
				if err := segmentQuery.Select("questions.id").Preload("CourseLink.Steps").Find(&allQuestions).Error; err != nil {
					return nil, &common.ErrorResponse{
						Code:    404,
						Message: "Failed to fetch questions",
						Details: err.Error(),
					}
				}

				for _, q := range allQuestions {
					if QuestionMeetsStepsRequirements(q, segmentData.Steps, *segmentData.StepsMode) {
						pickedQuestions = append(pickedQuestions, &TMPQ{
							ID:             q.ID,
							TimesUsed:      0,
							Answers:        convertAnswers(q.Answers),
							QuestionFormat: q.QuestionFormat,
						})
					}
				}
			case enums.CategoryFilterSQOR:
				// Add hand picked questions
				handpickedQuery := blockQuestionQuery.Session(&gorm.Session{})
				segmentQuestionGroups := make([]uint, len(segmentData.Questions))
				for q_ind, q := range segmentData.Questions {
					segmentQuestionGroups[q_ind] = q.ID
				}
				if err := handpickedQuery.Where("question_group_id IN (?)", segmentQuestionGroups).Find(&pickedQuestions).Error; err != nil {
					return nil, &common.ErrorResponse{
						Code:    404,
						Message: "Failed to fetch questions",
						Details: err.Error(),
					}
				}

				// All questions passing chapter&category&steps
				var allQuestions []models.Question
				if err := segmentQuery.Select("questions.id").Preload("CourseLink.Steps").Find(&allQuestions).Error; err != nil {
					return nil, &common.ErrorResponse{
						Code:    404,
						Message: "Failed to fetch questions",
						Details: err.Error(),
					}
				}

				for _, q := range allQuestions {
					if QuestionMeetsStepsRequirements(q, segmentData.Steps, *segmentData.StepsMode) && !slices.ContainsFunc(pickedQuestions, func(jt *TMPQ) bool {
						return jt.ID == q.ID
					}) {
						pickedQuestions = append(pickedQuestions, &TMPQ{
							ID:             q.ID,
							QuestionFormat: q.QuestionFormat,
							TimesUsed:      0,
							Answers:        convertAnswers(q.Answers),
						})
					}
				}
			default:
				panic(fmt.Sprintf("unexpected enums.CategoryFilterEnum: %#v", segmentData.FilterBy))
			}

			generatorCacheBlock.Segments = append(generatorCacheBlock.Segments, GeneratorCacheBlockSegment{
				ReqQuestionCount: segmentData.QuestionCount,
				QuestionPool:     pickedQuestions,
			})
		}
		generatorCache.Blocks = append(generatorCache.Blocks, generatorCacheBlock)
	}

	return &generatorCache, nil
}

func convertAnswers(answers []models.QuestionAnswer) []QuestionAnswer {
	newAnswers := make([]QuestionAnswer, len(answers))
	for i, answer := range answers {
		newAnswers[i] = QuestionAnswer{
			ID:         answer.ID,
			QuestionID: answer.QuestionID,
			AnswerID:   answer.AnswerID,
			Answer: &Answer{
				ID:          answer.Answer.ID,
				TimeToSolve: answer.Answer.TimeToSolve,
				Correct:     answer.Answer.Correct,
			},
		}
	}
	return newAnswers
}
