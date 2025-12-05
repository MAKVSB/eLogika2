package models

import (
	"time"

	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type Question struct {
	CommonModel
	ID          uint           `gorm:"primarykey"`
	CreatedAt   time.Time      ``
	CreatedByID uint           ``
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	UpdatedByID uint           ``
	DeletedAt   gorm.DeletedAt ``
	Version     uint           ``

	QuestionGroupID    uint                     ``                                           // Question origin tracking
	Title              string                   ``                                           // Title of the question
	Content            *TipTapContent           `gorm:"serializer:json;type:varbinary(max)"` // The text of the answer
	ContentFiles       []*File                  `gorm:"many2many:question_content_files;"`   // Files related to content
	TimeToRead         int                      ``                                           // Time to read the question
	TimeToProcess      int                      ``                                           // Time to common solution (building a graph or similar)
	QuestionType       enums.QuestionTypeEnum   ``                                           // Type of the question
	QuestionFormat     enums.QuestionFormatEnum ``                                           // Format of the question
	IncludeAnswerSpace bool                     ``                                           // Defines if a box of empty space should be included after open question
	ManagedBy          enums.CourseUserRoleEnum ``                                           // Role of user who manages it
	Active             bool                     ``                                           // If the question can be picked during test generation
	AnswerCount        uint                     ``

	QuestionGroup *QuestionGroup   ``
	Answers       []QuestionAnswer ``
	CheckedBy     []QuestionCheck  `gorm:"foreignKey:QuestionID"`
	CreatedBy     *User            ``
	UpdatedBy     *User            ``
	CourseLink    *CourseQuestion  ``
}

func (Question) TableName() string {
	return "questions"
}

func (Question) ApplyFilters(query *gorm.DB, filters []common.SearchRequestFilter, model any, extra map[string]any) (*gorm.DB, *common.ErrorResponse) {
	if filters != nil {
		// 1) Handle special cases and build a new slice without them
		var remainingFilters []common.SearchRequestFilter
		for _, filter := range filters {
			if filter.ID == "checkedBy" {
				if filter.Value == string(enums.QuestionCheckedByFilterChecked) {
					query = query.Where("EXISTS (SELECT 1 FROM question_checks WHERE question_checks.question_id = questions.id)")
				} else if filter.Value == string(enums.QuestionCheckedByFilterCheckedByMe) {
					query = query.Where("EXISTS (SELECT 1 FROM question_checks WHERE question_checks.question_id = questions.id AND question_checks.user_id = ?)", extra["userID"])
				} else if filter.Value == string(enums.QuestionCheckedByFilterUnchecked) {
					query = query.Where("NOT EXISTS (SELECT 1 FROM question_checks WHERE question_checks.question_id = questions.id)")
				} else if filter.Value == string(enums.QuestionCheckedByFilterUncheckedByMe) {
					query = query.Where("NOT EXISTS (SELECT 1 FROM question_checks WHERE question_checks.question_id = questions.id AND question_checks.user_id = ?)", extra["userID"])
				}
			} else {
				remainingFilters = append(remainingFilters, filter)
			}
		}

		query, err := CommonModel{}.ApplyFilters(query, remainingFilters, Question{}, nil, "")
		if err != nil {
			return query, err
		}
	}

	return query, nil
}
