package handlers

import (
	"slices"
	"time"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/course_item_terms/dtos"
	"elogika.vsb.cz/backend/modules/course_item_terms/helpers"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type StudentTermsListResponse struct {
	Items      []dtos.StudentTermDTO `json:"items"`
	ItemsCount int64                 `json:"itemsCount"`
}

// @Summary List all available terms of student
// @Tags Terms
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TermsListRequest true "Ability to filter results"
// @Success 200 {object} TermsListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/items/{courseItemId}/terms [get]
func ListForStudent(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		any,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// m := models.Term{}

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}
	// If not student
	if userRole != enums.CourseUserRoleStudent {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	// Get student studyForm
	var studentData models.CourseUser
	if err := initializers.DB.
		Model(&models.CourseUser{}).
		Where("user_id = ?", userData.ID).
		Where("course_id = ?", params.CourseID).
		First(&studentData).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch student tutors",
			Details: err.Error(),
		}
	}

	// Get all students classes
	var studentClassIds []uint
	if err := initializers.DB.
		Model(&models.ClassStudent{}).
		Where("user_id = ?", userData.ID).
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
		InnerJoins("INNER JOIN class_students on class_students.class_id = classes.id AND class_students.user_id = ? AND class_students.deleted_at is NULL", userData.ID).
		Pluck("class_tutors.id", &studentTutorIds).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch student tutors",
			Details: err.Error(),
		}
	}

	// Get all course items available for student
	var courseItemIds []uint
	if err := initializers.DB.
		Model(&models.CourseItem{}).
		Where("managed_by = ? OR (managed_by = ? AND created_by_id in ?)", enums.CourseUserRoleGarant, enums.CourseUserRoleTutor, studentTutorIds).
		Pluck("id", &courseItemIds).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch course items",
			Details: err.Error(),
		}
	}

	// Get all terms that are linked to available course items
	var terms []models.Term
	query := initializers.DB.
		Model(&models.Term{}).
		InnerJoins("CourseItem", initializers.DB.Where("study_form = ?", studentData.StudyForm)).
		Preload("CourseItem.Parent").
		Where("terms.course_id = ?", params.CourseID).
		Where("terms.course_item_id in ?", courseItemIds)

	if err := query.Find(&terms).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch terms",
			Details: err.Error(),
		}
	}

	// Get all terms that are already joined by student
	var allJoinedTerms []models.UserTerm
	if err := initializers.DB.Where("user_id = ?", userData.ID).Find(&allJoinedTerms).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to load all joined courses",
			Details: err.Error(),
		}
	}

	// Convert to DTOs
	dtoList := make([]dtos.StudentTermDTO, len(terms))
	for i, term := range terms {
		joinedStudentsCount, err := helpers.GetJoinedLocking(initializers.DB, term.ID, term.PerClass, studentClassIds, true)
		if err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to fetch numbers of joined students",
				Details: err,
			}
		}
		isJoined := slices.ContainsFunc(allJoinedTerms, func(jt models.UserTerm) bool {
			return jt.TermID == term.ID
		})
		canJoin := helpers.CanJoin(isJoined, term, joinedStudentsCount) == nil
		canLeave := helpers.CanLeave(isJoined, term) == nil
		willSignOut := false

		if canJoin {
			var alreadyJoinedTerm []models.UserTerm
			if err := initializers.DB.
				Where("user_id = ?", userData.ID).
				InnerJoins("Term", initializers.DB.Where("active_to > ?", time.Now())).
				InnerJoins("Term.CourseItem", initializers.DB.Where("Term__CourseItem.id = ?", term.CourseItemID)).
				Find(&alreadyJoinedTerm).Error; err != nil {
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to get already joined info",
					Details: err.Error(),
				}
			}

			if len(alreadyJoinedTerm) != 0 {
				willSignOut = true
			}

			for _, ajt := range alreadyJoinedTerm {
				if ajt.Term.SignOutTo.Before(time.Now()) {
					canJoin = false
				}
			}
		}

		dtoList[i] = dtos.StudentTermDTO{}.From(term, joinedStudentsCount, isJoined, canJoin, canLeave, willSignOut)
	}

	c.JSON(200, StudentTermsListResponse{
		Items:      dtoList,
		ItemsCount: int64(len(dtoList)),
	})

	return nil
}
