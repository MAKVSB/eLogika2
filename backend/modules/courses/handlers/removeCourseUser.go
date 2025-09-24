package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	authenums "elogika.vsb.cz/backend/modules/auth/enums"
	"elogika.vsb.cz/backend/modules/auth/tokens"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/utils"

	"github.com/gin-gonic/gin"
)

// @Description Request to insert user into class as tutor
type RemoveCourseUserRequest struct {
	UserID uint `json:"userId" binding:"required"`
}

// @Description Updated list of students
type RemoveCourseUserResponse struct {
	Success bool `json:"success"`
}

// @Summary Add user to course
// @Tags Courses
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body RemoveCourseUserRequest true "New data for class"
// @Success 200 {object} RemoveCourseUserResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/users [delete]
func RemoveCourseUser(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		RemoveCourseUserRequest,
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
	if userRole != enums.CourseUserRoleAdmin {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	var courseUser models.CourseUser
	if err := initializers.DB.
		Where("user_id = ?", reqData.UserID).
		Where("course_id = ?", params.CourseID).
		First(&courseUser).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch course user",
		}
	}

	if err := initializers.DB.
		Where("user_id = ?", reqData.UserID).
		Where("course_id = ?", params.CourseID).
		Delete(&models.CourseUser{}).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to remove course user",
			Details: err.Error(),
		}
	}

	// TODO TODO TODO invalidate user token so i can make sure that user will loose access

	accessToken := tokens.AccessToken{}
	accessToken.UserID = reqData.UserID
	accessToken.TokenType = authenums.JWTTokenTypeAccess
	accessToken.InvalidateByUser()

	c.JSON(200, RemoveCourseUserResponse{
		Success: true,
	})

	return nil
}
