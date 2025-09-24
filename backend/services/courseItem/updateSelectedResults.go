package services_course_item

import (
	"fmt"
	"slices"
	"time"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type TempTermResult struct {
	ResultID uint
	Points   float64
}

type TempTerm struct {
	TermID        uint
	GroupTerm     bool //If term is for group
	TermActiveEnd time.Time
	Results       map[uint]TempTermResult //courseItemID
}

func UpdateSelectedResultsRoot(dbRef *gorm.DB, courseId uint, studentId uint, studyForm enums.StudyFormEnum) *common.ErrorResponse {
	var results []*models.CourseItemResult
	if err := dbRef.
		Where("student_id = ?", studentId).
		Order("created_at DESC").
		Find(&results).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch student restuls",
			Details: err.Error(),
		}
	}

	var rootCourseItems []models.CourseItem
	if err := dbRef.
		Where("course_id = ?", courseId).
		Where("parent_id is NULL").
		Where("study_form = ?", studyForm).
		Joins("ActivityDetail").
		Joins("GroupDetail").
		Joins("TestDetail").
		Preload("Terms", func(db *gorm.DB) *gorm.DB {
			return db.
				InnerJoins("INNER JOIN user_terms on user_terms.term_id = Terms.id AND user_terms.user_id = ?", studentId).
				Order("user_terms.created_at DESC")
		}).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("ActivityDetail").
				Joins("GroupDetail").
				Joins("TestDetail").
				Preload("Terms", func(db *gorm.DB) *gorm.DB {
					return db.
						InnerJoins("INNER JOIN user_terms on user_terms.term_id = Terms.id AND user_terms.user_id = ?", studentId).
						Order("user_terms.created_at DESC")
				})
		}).
		Find(&rootCourseItems).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch course items",
			Details: err.Error(),
		}
	}

	for _, rootCourseItem := range rootCourseItems {
		err := UpdateCourseItemSelectedResult(dbRef, &rootCourseItem, &results, studentId)
		if err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to update active results",
			}
		}
	}

	return nil
}

func UpdateSelectedResults(dbRef *gorm.DB, courseId uint, rootCourseItemId uint, studentId uint) *common.ErrorResponse {
	var courseUser *models.CourseUser
	if err := dbRef.
		Where("user_id = ?", studentId).
		First(&courseUser).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch student restuls",
			Details: err.Error(),
		}
	}

	var results []*models.CourseItemResult
	if err := dbRef.
		Where("student_id = ?", studentId).
		Order("created_at DESC").
		Find(&results).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch student restuls",
			Details: err.Error(),
		}
	}

	var rootCourseItems []models.CourseItem
	if err := dbRef.
		Where("course_id = ?", courseId).
		Where("parent_id is NULL").
		Where("study_form = ?", courseUser.StudyForm).
		Joins("ActivityDetail").
		Joins("GroupDetail").
		Joins("TestDetail").
		Preload("Terms", func(db *gorm.DB) *gorm.DB {
			return db.
				InnerJoins("INNER JOIN user_terms on user_terms.term_id = Terms.id AND user_terms.user_id = ?", studentId).
				Order("user_terms.created_at DESC")
		}).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("ActivityDetail").
				Joins("GroupDetail").
				Joins("TestDetail").
				Preload("Terms", func(db *gorm.DB) *gorm.DB {
					return db.
						InnerJoins("INNER JOIN user_terms on user_terms.term_id = Terms.id AND user_terms.user_id = ?", studentId).
						Order("user_terms.created_at DESC")
				})
		}).
		First(&rootCourseItems, rootCourseItemId).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch course items",
			Details: err.Error(),
		}
	}

	for _, rootCourseItem := range rootCourseItems {
		err := UpdateCourseItemSelectedResult(dbRef, &rootCourseItem, &results, studentId)
		if err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to update active results",
			}
		}
	}

	return nil
}

func UpdateCourseItemSelectedResult(dbRef *gorm.DB, courseItem *models.CourseItem, results *[]*models.CourseItemResult, studentId uint) *common.ErrorResponse {

	courseItemIDs := []uint{}

	courseItemIDs = append(courseItemIDs, courseItem.ID)
	for _, courseItemChildren := range courseItem.Children {
		courseItemIDs = append(courseItemIDs, courseItemChildren.ID)
	}

	if err := dbRef.
		Model(&models.CourseItemResult{}).
		Where("course_item_id in ?", courseItemIDs).
		Where("student_id = ?", studentId).
		Update("selected", false).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to unselect results",
			Details: err.Error(),
		}
	}

	switch courseItem.Type {
	case enums.CourseItemTypeActivity, enums.CourseItemTypeTest:
		tempResult := EvaluateByAttempt(results, courseItem, studentId, nil)
		if tempResult != nil {
			if err := dbRef.
				Model(&models.CourseItemResult{}).
				Where("id = ?", tempResult.ResultID).
				Update("selected", true).Error; err != nil {
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to unselect results",
					Details: err.Error(),
				}
			}
		}
	case enums.CourseItemTypeGroup:
		mustBeInSameTerm := false
		if len(courseItem.Terms) != 0 {
			mustBeInSameTerm = true
		}

		if !mustBeInSameTerm {
			// Group is not set up in a way that requires selected items to be in the same term
			for _, courseItemChildren := range courseItem.Children {
				tempResult := EvaluateByAttempt(results, courseItemChildren, studentId, nil)
				if tempResult != nil {
					if err := dbRef.
						Model(&models.CourseItemResult{}).
						Where("id = ?", tempResult.ResultID).
						Update("selected", true).Error; err != nil {
						return &common.ErrorResponse{
							Code:    500,
							Message: "Failed to unselect results",
							Details: err.Error(),
						}
					}
				}
			}
		} else {
			// Must parse to terms, sort and the pick the correct one
			terms := []*TempTerm{}

			// Load all relevant MAIN terms into memory
			switch courseItem.EvaluateByAttempt {
			case enums.EvaluateByAttemptBest:
				for _, courseItemTerm := range courseItem.Terms {
					tempTerm := TempTerm{
						TermID:        courseItemTerm.ID,
						GroupTerm:     true,
						TermActiveEnd: courseItemTerm.ActiveFrom,
						Results:       map[uint]TempTermResult{},
					}

					for _, courseItemChildren := range courseItem.Children {
						tempResult := EvaluateByAttempt(results, courseItemChildren, studentId, &courseItemTerm.ID)
						if tempResult != nil {
							tempTerm.Results[courseItemChildren.ID] = *tempResult
						}
					}
					terms = append(terms, &tempTerm)
				}
			case enums.EvaluateByAttemptLast:
				courseItemTerm := courseItem.Terms[0]

				tempTerm := TempTerm{
					TermID:        courseItemTerm.ID,
					GroupTerm:     true,
					TermActiveEnd: courseItemTerm.ActiveFrom,
					Results:       map[uint]TempTermResult{},
				}

				for _, courseItemChildren := range courseItem.Children {
					tempResult := EvaluateByAttempt(results, courseItemChildren, studentId, &courseItemTerm.ID)
					if tempResult != nil {
						tempTerm.Results[courseItemChildren.ID] = *tempResult
					}
				}
				terms = append(terms, &tempTerm)
			default:
				panic(fmt.Sprintf("unexpected enums.EvaluateByAttemptEnum: %#v", courseItem.EvaluateByAttempt))
			}

			// Get child created terms
			for _, courseItemChildren := range courseItem.Children {
				switch courseItem.EvaluateByAttempt {

				case enums.EvaluateByAttemptLast:
					if len(courseItemChildren.Terms) != 0 {
						courseItemChildrenTerm := courseItemChildren.Terms[0]
						tempTerm := TempTerm{
							TermID:        courseItemChildrenTerm.ID,
							GroupTerm:     false,
							TermActiveEnd: courseItemChildrenTerm.ActiveFrom,
							Results:       map[uint]TempTermResult{},
						}

						tempResult := EvaluateByAttempt(results, courseItemChildren, studentId, &courseItemChildrenTerm.ID)
						if tempResult != nil {
							tempTerm.Results[courseItemChildren.ID] = *tempResult
						}
						terms = append(terms, &tempTerm)
					}
				case enums.EvaluateByAttemptBest:
					for _, courseItemChildrenTerm := range courseItemChildren.Terms {
						tempTerm := TempTerm{
							TermID:        courseItemChildrenTerm.ID,
							GroupTerm:     false,
							TermActiveEnd: courseItemChildrenTerm.ActiveFrom,
							Results:       map[uint]TempTermResult{},
						}

						tempResult := EvaluateByAttempt(results, courseItemChildren, studentId, &courseItemChildrenTerm.ID)
						if tempResult != nil {
							tempTerm.Results[courseItemChildren.ID] = *tempResult
						}
						terms = append(terms, &tempTerm)
					}
				default:
					panic(fmt.Sprintf("unexpected enums.EvaluateByAttemptEnum: %#v", courseItem.EvaluateByAttempt))
				}
			}

			// Sort from oldest to newest
			slices.SortFunc(terms, func(a, b *TempTerm) int {
				return a.TermActiveEnd.Compare(b.TermActiveEnd)
			})

			// Choose the best/newest
			var finalTempTerm *TempTerm

			switch courseItem.EvaluateByAttempt {
			case enums.EvaluateByAttemptLast:
				for _, v := range terms {
					if finalTempTerm == nil {
						if v.GroupTerm {
							finalTempTerm = v
						}
					} else {
						if v.GroupTerm {
							finalTempTerm = v
						} else {
							for v_r_id, v_r := range v.Results {
								finalTempTerm.Results[v_r_id] = v_r
							}
						}
					}
				}
			case enums.EvaluateByAttemptBest:
				// Combine partial terms into main terms
				var completeTempTerms []*TempTerm
				var completeTempTermLast *TempTerm

				for _, v := range terms {
					if completeTempTermLast == nil {
						if v.GroupTerm {
							completeTempTerms = append(completeTempTerms, v)
							completeTempTermLast = v
						}
					} else {
						if v.GroupTerm {
							completeTempTerms = append(completeTempTerms, v)
							completeTempTermLast = v
						} else {
							for v_r_id, v_r := range v.Results {
								val, ok := completeTempTermLast.Results[v_r_id]
								if !ok || val.Points < v_r.Points {
									completeTempTermLast.Results[v_r_id] = v_r
								}
							}
						}
					}
				}

				// Pick the best effort
				bestPoints := float64(0)
				for _, v := range completeTempTerms {
					pointsSum := float64(0)

					for _, v_r := range v.Results {
						pointsSum += v_r.Points
					}

					if pointsSum >= bestPoints {
						finalTempTerm = v
						bestPoints = pointsSum
					}
				}
			default:
				panic(fmt.Sprintf("unexpected enums.EvaluateByAttemptEnum: %#v", courseItem.EvaluateByAttempt))
			}

			if finalTempTerm != nil {
				for _, v := range finalTempTerm.Results {
					if err := dbRef.
						Model(&models.CourseItemResult{}).
						Where("id = ?", v.ResultID).
						Update("selected", true).Error; err != nil {
						return &common.ErrorResponse{
							Code:    500,
							Message: "Failed to unselect results",
							Details: err.Error(),
						}
					}
				}
			}
		}
	default:
		panic(fmt.Sprintf("unexpected enums.CourseItemTypeEnum: %#v", courseItem.Type))
	}
	return nil
}

func FindInResults(results *[]*models.CourseItemResult, courseItemID uint, studentId uint, termId *uint) []*models.CourseItemResult {
	var result []*models.CourseItemResult
	for _, v := range *results {
		if v.CourseItemID == courseItemID && v.StudentID == studentId {
			if termId == nil || v.TermID == *termId {
				result = append(result, v)
			}
		}
	}
	return result
}

func EvaluateByAttemptBest(results *[]*models.CourseItemResult, courseItemID uint, studentId uint, termId *uint) (uint, float64) {
	var bestResID uint
	var bestResPoints float64

	item_results := FindInResults(results, courseItemID, studentId, termId)
	for _, v := range item_results {
		if v.Points > bestResPoints {
			bestResPoints = v.Points
			bestResID = v.ID
		}
	}

	return bestResID, bestResPoints
}

func EvaluateByAttempt(results *[]*models.CourseItemResult, courseItem *models.CourseItem, studentId uint, termId *uint) *TempTermResult {
	switch courseItem.EvaluateByAttempt {
	case enums.EvaluateByAttemptBest:
		bestResID, bestResPoints := EvaluateByAttemptBest(results, courseItem.ID, studentId, termId)

		if bestResID != 0 {
			return &TempTermResult{
				ResultID: bestResID,
				Points:   bestResPoints,
			}
		}
	case enums.EvaluateByAttemptLast:
		item_results := FindInResults(results, courseItem.ID, studentId, termId)
		if len(item_results) != 0 {
			return &TempTermResult{
				ResultID: item_results[0].ID,
				Points:   item_results[0].Points,
			}
		}
	default:
		panic(fmt.Sprintf("unexpected enums.EvaluateByAttemptEnum: %#v", courseItem.EvaluateByAttempt))
	}
	return nil
}
