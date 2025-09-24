package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/classes/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Request to insert user into class as tutor
type RemoveStudentRequest struct {
	UserID uint `json:"userId" binding:"required"`
}

// @Description Updated list of tutors
type RemoveStudentResponse struct {
	Students []dtos.ClassUserDTO `json:"students"`
}

// @Summary Removes tutor from class
// @Tags Classes
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body RemoveStudentRequest true "New data for class"
// @Success 200 {object} RemoveStudentResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/ [post] // TODO
func RemoveStudent(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
			ClassID  uint `uri:"classId"  binding:"required"`
		},
		RemoveStudentRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		return err
	}

	classRepo := repositories.NewClassRepository()
	// If not admin, garant or tutor
	if userRole == enums.CourseUserRoleAdmin || userRole == enums.CourseUserRoleGarant {
		// Permission to update every class
	} else if userRole == enums.CourseUserRoleTutor {
		// Can only update his own class
		_, err = classRepo.GetClassByIDTutor(initializers.DB, params.CourseID, params.ClassID, userData.ID, false, nil)
	} else {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
	if err != nil {
		return err
	}

	var classStudent *models.ClassStudent
	if err := initializers.DB.
		Where("user_id = ?", reqData.UserID).
		Where("class_id = ?", params.ClassID).
		Delete(&classStudent).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to remove student from class",
		}
	}

	var classStudents []models.ClassStudent
	if err := initializers.DB.
		Where("class_id = ?", params.ClassID).
		InnerJoins("User").
		Find(&classStudents).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch class students",
			Details: err.Error(),
		}
	}

	// Convert to DTOs
	dtoList := make([]dtos.ClassUserDTO, len(classStudents))
	for i, c := range classStudents {
		dtoList[i] = dtos.ClassUserDTO{}.From(c.User)
	}

	c.JSON(200, RemoveStudentResponse{
		Students: dtoList,
	})
	return nil
}
