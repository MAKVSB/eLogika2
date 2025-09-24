package helpers

import (
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
)

func GetTestAttemptsLeft(courseItemChildId uint, termId uint, userId uint) (uint, *common.ErrorResponse) {
	var courseItemGroup *models.CourseItem
	var courseItemChild models.CourseItem
	var term models.Term

	if err := initializers.DB.
		First(&courseItemChild, courseItemChildId).Error; err != nil {
		return 0, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to load course item data",
		}
	}

	if courseItemChild.ParentID != nil {
		if err := initializers.DB.
			First(&courseItemGroup, courseItemChild.ParentID).Error; err != nil {
			return 0, &common.ErrorResponse{
				Code:    500,
				Message: "Failed to load course item group data",
			}
		}
	}

	if err := initializers.DB.
		First(&term, termId).Error; err != nil {
		return 0, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to load term data",
		}
	}

	// Actual calculation
	courseItemMaxAttempts := int64(courseItemChild.MaxAttempts)
	if courseItemGroup != nil {
		if courseItemGroup.MaxAttempts < courseItemChild.MaxAttempts {
			courseItemMaxAttempts = int64(courseItemGroup.MaxAttempts)
		}
	}

	// How many instances have been run in this courseItem
	var courseItemAttempts int64
	if courseItemChild.ParentID != nil {
		if err := initializers.DB.
			Model(models.TestInstance{}).
			Where("participant_id = ?", userId).
			Where("course_item_id = ?", courseItemChild.ID).
			Count(&courseItemAttempts).Error; err != nil {
			return 0, &common.ErrorResponse{
				Code:    500,
				Message: "Failed to load course item instances",
			}
		}
	}

	// How many instances have been run in this term
	var termAttempts int64
	if courseItemChild.ParentID != nil {
		if err := initializers.DB.
			Model(models.TestInstance{}).
			Where("participant_id = ?", userId).
			Where("term_id = ?", termId).
			Where("course_item_id = ?", courseItemChild.ID).
			Count(&termAttempts).Error; err != nil {
			return 0, &common.ErrorResponse{
				Code:    500,
				Message: "Failed to load term instances",
			}
		}
	}

	z_test := uint(max(0, courseItemMaxAttempts-courseItemAttempts)) // How many attempts is left for test
	z_term := uint(max(0, int64(term.Tries)-termAttempts))           // How many attempts is left for term

	return min(z_test, z_term), nil
}
