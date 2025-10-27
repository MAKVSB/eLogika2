package handlers

import (
	"time"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/course_item_terms/helpers"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Newly created question
type TermsJoinResponse struct {
	Success bool `json:"success"`
}

type TermsJoinRequest struct {
	UserID *uint `json:"userId"`
}

// @Summary Join term as student
// @Tags Terms
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param termId path int true "ID of the corresponding term"
// @Success 200 {object} TermsJoinResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/terms/{termId} [post]
func UserJoin(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
			TermID       uint `uri:"termId" binding:"required"`
		},
		TermsJoinRequest,
	](c)
	if err != nil {
		return err
	}

	// // TODO validate from here
	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}

	userId := userData.ID
	// If not admin, garant, or tutor
	if userRole == enums.CourseUserRoleAdmin {
		if reqData.UserID != nil {
			userId = *reqData.UserID
		} else {
			return &common.ErrorResponse{
				Code:    422,
				Message: "Student not specified",
			}
		}
	} else if userRole == enums.CourseUserRoleGarant {
		if reqData.UserID != nil {
			userId = *reqData.UserID
		} else {
			return &common.ErrorResponse{
				Code:    422,
				Message: "Student not specified",
			}
		}

		var courseItem models.CourseItem
		if err := initializers.DB.
			Where("managed_by = ?", enums.CourseUserRoleGarant).
			Find(&courseItem, params.CourseItemID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permission for this item",
			}
		}
	} else if userRole == enums.CourseUserRoleTutor {
		if reqData.UserID != nil {
			userId = *reqData.UserID
		} else {
			return &common.ErrorResponse{
				Code:    422,
				Message: "Student not specified",
			}
		}

		var courseItem models.CourseItem
		if err := initializers.DB.
			Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userId).
			Find(&courseItem, params.CourseItemID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permission for this item",
			}
		}
	} else if userRole == enums.CourseUserRoleStudent {
		// Get all students classes
		var studentClassIds []uint
		if err := initializers.DB.
			Model(&models.ClassStudent{}).
			Where("user_id = ?", userId).
			Pluck("class_id", &studentClassIds).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to fetch student tutors",
				Details: err.Error(),
			}
		}

		// Get all tutors that teach the student
		var studentTutorIds []uint
		if err := initializers.DB.
			Model(&models.ClassTutor{}).
			InnerJoins("INNER JOIN classes on classes.id = class_tutors.class_id AND classes.type = ?", enums.ClassTypeC).
			InnerJoins("INNER JOIN class_students on class_students.class_id = classes.id AND class_students.user_id = ? AND class_students.deleted_at is NULL", userId).
			Pluck("class_tutors.id", &studentTutorIds).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to fetch student tutors",
				Details: err.Error(),
			}
		}

		// Get all course items available for student
		var courseItem []models.CourseItem
		if err := initializers.DB.
			Model(&models.CourseItem{}).
			Where("managed_by = ? OR (managed_by = ? AND created_by_id in ?)", enums.CourseUserRoleGarant, enums.CourseUserRoleTutor, studentTutorIds).
			Find(&courseItem, params.CourseItemID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to fetch course items",
				Details: err.Error(),
			}
		}
		// TODO Check if user can modify courseItem
	} else {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	// Check if user in course
	var courseUser models.CourseUser
	if err := initializers.DB.
		Where("course_id = ?", params.CourseID).
		Where("user_id = ?", userId).
		Where("roles like '%?%'", enums.CourseUserRoleStudent).
		Find(&courseUser, params.CourseItemID).Error; err != nil {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permission for this item",
		}
	}

	transaction := initializers.DB.Begin()

	// Get term data
	var term models.Term
	if err := transaction.
		Where("id = ?", params.TermID).
		First(&term).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to get term info",
			Details: err.Error(),
		}
	}

	// Get all students classes
	var studentClassIds []uint
	if err := transaction.
		Model(&models.ClassStudent{}).
		Where("user_id = ?", userId).
		Pluck("class_id", &studentClassIds).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch student tutors",
			Details: err.Error(),
		}
	}

	// Get all tutors that teach the student
	var studentTutorIds []uint
	if err := transaction.
		Model(&models.ClassTutor{}).
		InnerJoins("INNER JOIN classes on classes.id = class_tutors.class_id AND classes.type = ?", enums.ClassTypeC).
		InnerJoins("INNER JOIN class_students on class_students.class_id = classes.id AND class_students.user_id = ? AND class_students.deleted_at is NULL", userId).
		Pluck("class_tutors.id", &studentTutorIds).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch student tutors",
			Details: err.Error(),
		}
	}

	joinedStudentsCount, err := helpers.GetJoinedLocking(transaction, params.TermID, term.PerClass, studentClassIds, true)
	if err != nil {
		transaction.Rollback()
		return err
	}

	var studentJoinedCheck models.UserTerm
	if err := transaction.
		Model(&models.UserTerm{}).
		Where("user_id = ?", userId).
		Where("term_id = ?", params.TermID).
		Find(&studentJoinedCheck).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch student tutors",
			Details: err.Error(),
		}
	}

	isJoined := false
	if studentJoinedCheck.ID != 0 {
		isJoined = true
	}

	canJoinError := helpers.CanJoin(isJoined, term, joinedStudentsCount)
	if canJoinError != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to sign in to term",
			Details: canJoinError.Error(),
		}
	}

	var alreadyJoinedTerm []models.UserTerm
	if err := transaction.
		Where("user_id = ?", userId).
		InnerJoins("Term", initializers.DB.Where("active_to > ?", time.Now())).
		InnerJoins("Term.CourseItem", initializers.DB.Where("Term__CourseItem.id = ?", term.CourseItemID)).
		Find(&alreadyJoinedTerm).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to get already joined info",
			Details: err.Error(),
		}
	}

	for _, ajt := range alreadyJoinedTerm {
		if ajt.Term.SignOutTo.Before(time.Now()) {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    409,
				Message: "Cannot sign out from joined term",
			}
		}

		ajt.DeletedByID = &userData.ID

		if err := transaction.
			Save(&ajt).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to sign out from joined term",
			}
		}

		if err := transaction.
			Delete(&ajt).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to sign out from joined term",
			}
		}
	}

	termJoin := models.UserTerm{
		ID:          0,
		CreatedByID: userData.ID,
		UserID:      userId,
		TermID:      params.TermID,
	}

	if err := transaction.Save(&termJoin).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to join term",
			Details: err.Error(),
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	c.JSON(200, TermsJoinResponse{
		Success: true,
	})

	return nil
}
