package handlers

import (
	"time"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/course_item_terms/dtos"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Request to update term
type TermsUpdateRequest struct {
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

	Version uint `json:"version"` // Version signature to prevent concurrency problems
}

// @Description Updated term
type TermsUpdateResponse struct {
	Data dtos.TermDTO `json:"data"`
}

// @Summary Update term
// @Tags CourseItems
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TermsUpdateRequest true "New data for question"
// @Success 200 {object} TermsUpdateResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/items/{courseItemId}/terms/{termId} [put]
func Update(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
			TermID       uint `uri:"termId" binding:"required"`
		},
		TermsUpdateRequest,
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

	termService := services.TermService{}
	term, err := termService.GetTermByID(initializers.DB, params.CourseID, params.CourseItemID, params.TermID, userData.ID, userRole, nil, true, &reqData.Version)
	if err != nil {
		return err
	}

	term.Version = term.Version + 1
	term.Name = reqData.Name
	term.ActiveFrom = reqData.ActiveFrom
	term.ActiveTo = reqData.ActiveTo
	term.RequiresSign = reqData.RequiresSign
	term.SignInFrom = reqData.SignInFrom
	term.SignInTo = reqData.SignInTo
	term.SignOutFrom = reqData.SignOutFrom
	term.SignOutTo = reqData.SignOutTo
	// term.OfflineTo =     reqData.OfflineTo
	term.Classroom = reqData.Classroom
	term.StudentsMax = reqData.StudentsMax
	term.Tries = reqData.Tries

	if err := initializers.DB.Save(&term).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update term",
			Details: err.Error(),
		}
	}

	term, err = termService.GetTermByID(initializers.DB, params.CourseID, params.CourseItemID, params.TermID, userData.ID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	c.JSON(200, TermsUpdateResponse{
		Data: dtos.TermDTO{}.From(term),
	})
	return nil
}
