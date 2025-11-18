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
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Request to insert user into class as tutor
type AddStudentRequest struct {
	UserID uint `json:"userId" binding:"required"`
}

// @Description Updated list of students
type AddStudentResponse struct {
	Students []dtos.ClassUserDTO `json:"students"`
}

// @Summary Add student to class
// @Tags Classes
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body AddStudentRequest true "New data for class"
// @Success 200 {object} AddStudentResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/classes/{classId}/students [post]
func AddStudent(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
			ClassID  uint `uri:"classId" binding:"required"`
		},
		AddStudentRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}

	classService := services.NewClassService(repositories.NewClassRepository())
	_, err = classService.GetClassByID(initializers.DB, params.CourseID, params.ClassID, userData.ID, userRole, nil, false, nil)
	if err != nil {
		return err
	}

	var courseUser models.CourseUser
	if err := initializers.DB.
		Where("user_id = ?", reqData.UserID).
		Where("course_id = ?", params.CourseID).
		First(&courseUser).Error; err != nil {
		// TODO is specific error = "User not in course"
		return &common.ErrorResponse{
			Code:    404,
			Message: "User not in course",
		}
	}
	if courseUser.NotHasRole(enums.CourseUserRoleStudent) {
		return &common.ErrorResponse{
			Code:    409,
			Message: "User does not have student role",
		}
	}

	var classStudent *models.ClassStudent
	if err := initializers.DB.
		Where("user_id = ?", reqData.UserID).
		Where("class_id = ?", params.ClassID).
		Find(&classStudent).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to validate if user already is student of specific class",
			Details: err.Error(),
		}
	}
	if classStudent != nil && classStudent.ID != 0 {
		return &common.ErrorResponse{
			Code:    409,
			Message: "User already is student for this class",
		}
	}

	classStudent = &models.ClassStudent{
		UserID:  reqData.UserID,
		ClassID: params.ClassID,
	}

	if err := initializers.DB.Save(&classStudent).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to insert class student",
		}
	}

	var classStudents []models.ClassStudent
	if err := initializers.DB.
		Where("class_id = ?", params.ClassID).
		InnerJoins("User").
		Find(&classStudents).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch class student",
		}
	}

	// Convert to DTOs
	dtoList := make([]dtos.ClassUserDTO, len(classStudents))
	for i, c := range classStudents {
		dtoList[i] = dtos.ClassUserDTO{}.From(c.User)
	}

	c.JSON(200, AddStudentResponse{
		Students: dtoList,
	})
	return nil
}
