package handlers

import (
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/classes/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Request to insert user into class as tutor
type AddTutorRequest struct {
	UserID uint `json:"userId" binding:"required"`
}

// @Description Updated list of tutors
type AddTutorResponse struct {
	Tutors []dtos.ClassUserDTO `json:"tutors"`
}

// @Summary Add tutor to class
// @Tags Classes
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body AddTutorRequest true "New data for class"
// @Success 200 {object} AddTutorResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/classes/{classId}/tutors [post]
func AddTutor(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
			ClassID  uint `uri:"classId" binding:"required"`
		},
		AddTutorRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// If not admin or garant
	if userRole != enums.CourseUserRoleAdmin && userRole != enums.CourseUserRoleGarant {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	transaction := initializers.DB.Begin()

	var courseUser models.CourseUser
	if err := transaction.
		Where("user_id = ?", reqData.UserID).
		Where("course_id = ?", params.CourseID).
		First(&courseUser).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    404,
			Message: "User not in course",
		}
	}

	if courseUser.NotHasRole(enums.CourseUserRoleTutor) && courseUser.NotHasRole(enums.CourseUserRoleGarant) {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    409,
			Message: "User does not have tutor or garant role",
		}
	}

	var classTutor *models.ClassTutor
	if err := transaction.
		Where("user_id = ?", reqData.UserID).
		Where("class_id = ?", params.ClassID).
		Find(&classTutor).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to validate if user already is tutor",
		}
	}
	if classTutor != nil && classTutor.ID != 0 {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    409,
			Message: "User already is tutor for this class",
		}
	}

	classTutor = &models.ClassTutor{
		UserID:  reqData.UserID,
		ClassID: params.ClassID,
	}

	if err := transaction.Save(&classTutor).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to insert class tutor",
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	var classTutors []models.ClassTutor
	if err := initializers.DB.
		Where("class_id = ?", params.ClassID).
		InnerJoins("User").
		Find(&classTutors).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch class tutors",
		}
	}

	// Convert to DTOs
	dtoList := make([]dtos.ClassUserDTO, len(classTutors))
	for i, c := range classTutors {
		dtoList[i] = dtos.ClassUserDTO{}.From(c.User)
	}

	c.JSON(200, AddTutorResponse{
		Tutors: dtoList,
	})
	return nil
}
