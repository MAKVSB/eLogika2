package handlers

import (
	"time"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/course_item_terms/dtos"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Request to insert new question
type TermsInsertRequest struct {
	Name         string    `json:"name"`
	ActiveFrom   time.Time `json:"activeFrom"`
	ActiveTo     time.Time `json:"activeTo"`
	RequiresSign bool      `json:"requiresSign"`
	SignInFrom   time.Time `json:"signInFrom"`
	SignInTo     time.Time `json:"signInTo"`
	SignOutFrom  time.Time `json:"signOutFrom"`
	SignOutTo    time.Time `json:"signOutTo"`
	// OfflineTo:    n,
	Classroom   string `json:"classroom"`
	StudentsMax uint   `json:"studentsMax"`
	Tries       uint   `json:"tries"`
}

// @Description Newly created question
type TermsInsertResponse struct {
	Data dtos.TermDTO `json:"data"`
}

// @Summary Create new term
// @Tags Terms
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TermsInsertRequest true "New data for question"
// @Success 200 {object} TermsInsertResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/items/{courseItemId}/terms [post]
func Insert(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
		},
		TermsInsertRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		return err
	}
	// If not admin, garant, or tutor
	if userRole != enums.CourseUserRoleAdmin && userRole != enums.CourseUserRoleGarant && userRole != enums.CourseUserRoleTutor {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	// Can user see the course item
	cis := services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
	courseItem, err := cis.GetCourseItemByID(initializers.DB, params.CourseID, params.CourseItemID, userData.ID, userRole, nil, false, nil)
	if err != nil {
		return err
	}
	if !courseItem.Editable {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	term := &models.Term{
		ID:           0,
		Version:      1,
		CourseID:     params.CourseID,
		CourseItemID: params.CourseItemID,

		Name:         reqData.Name,
		ActiveFrom:   reqData.ActiveFrom,
		ActiveTo:     reqData.ActiveTo,
		RequiresSign: reqData.RequiresSign,
		SignInFrom:   reqData.SignInFrom,
		SignInTo:     reqData.SignInTo,
		SignOutFrom:  reqData.SignOutFrom,
		SignOutTo:    reqData.SignOutTo,
		// OfflineTo:    reqData.OfflineTo,
		Classroom:   reqData.Classroom,
		StudentsMax: reqData.StudentsMax,
		Tries:       reqData.Tries,
	}

	if err := initializers.DB.Save(&term).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to insert term",
			Details: err.Error(),
		}
	}

	termService := services.TermService{}
	term, err = termService.GetTermByID(initializers.DB, params.CourseID, params.CourseItemID, term.ID, userData.ID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	c.JSON(200, TermsInsertResponse{
		Data: dtos.TermDTO{}.From(term),
	})
	return nil
}
