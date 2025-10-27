package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Request to insert user into class as tutor
type AddCourseUserRequest struct {
	UserID    uint                     `json:"userId" binding:"required"`
	Role      enums.CourseUserRoleEnum `json:"role" binding:"required"`
	StudyForm *enums.StudyFormEnum     `json:"studyForm"`
}

// @Description Updated list of students
type AddCourseUserResponse struct {
	Success bool `json:"success"`
}

// @Summary Add student to course
// @Tags Courses
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body AddCourseUserRequest true "New data for class"
// @Success 200 {object} AddCourseUserResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/users [post]
func AddCourseUser(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		AddCourseUserRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}
	// If not admin, garant, or tutor
	if userRole != enums.CourseUserRoleAdmin && userRole != enums.CourseUserRoleGarant {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	transaction := initializers.DB.Begin()

	var courseUser *models.CourseUser
	if err := transaction.
		Where("user_id = ?", reqData.UserID).
		Where("course_id = ?", params.CourseID).
		Find(&courseUser).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch course user",
		}
	}

	if courseUser == nil || courseUser.ID == 0 {
		courseUser = &models.CourseUser{
			CourseID: params.CourseID,
			UserID:   reqData.UserID,
		}
	}

	if reqData.Role == enums.CourseUserRoleStudent {
		if reqData.StudyForm == nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    409,
				Message: "Student must have studyForm attribute set",
				Details: "User already has role",
			}
		}
		courseUser.StudyForm = reqData.StudyForm
	}

	alreadyHasRole := false
	for _, v := range courseUser.Roles {
		if v == reqData.Role {
			alreadyHasRole = true
		}
	}

	if !alreadyHasRole {
		courseUser.Roles = append(courseUser.Roles, reqData.Role)
	}

	utils.DebugPrintJSON(courseUser)

	if err := transaction.Save(&courseUser).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to insert course user",
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

	c.JSON(200, AddCourseUserResponse{
		Success: true,
	})
	return nil
}
