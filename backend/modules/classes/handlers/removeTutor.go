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
type RemoveTutorRequest struct {
	UserID uint `json:"userId" binding:"required"`
}

// @Description Updated list of tutors
type RemoveTutorResponse struct {
	Tutors []dtos.ClassUserDTO `json:"tutors"`
}

// @Summary Removes tutor from class
// @Tags Classes
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body RemoveTutorRequest true "New data for class"
// @Success 200 {object} RemoveTutorResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/classes/{classId}/tutors [delete]
func RemoveTutor(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
			ClassID  uint `uri:"classId"  binding:"required"`
		},
		RemoveTutorRequest,
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

	var classTutor *models.ClassTutor
	if err := initializers.DB.
		Where("user_id = ?", reqData.UserID).
		Where("class_id = ?", params.ClassID).
		Delete(&classTutor).Error; err != nil {
		return &common.ErrorResponse{
			Code:    404,
			Message: "User is not a tutor",
		}
	}

	var classTutors []models.ClassTutor
	if err := initializers.DB.
		Where("class_id = ?", params.ClassID).
		InnerJoins("User").
		Find(&classTutors).Error; err != nil {
		return &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch class tutors",
		}
	}

	// Convert to DTOs
	dtoList := make([]dtos.ClassUserDTO, len(classTutors))
	for i, c := range classTutors {
		dtoList[i] = dtos.ClassUserDTO{}.From(c.User)
	}

	c.JSON(200, RemoveTutorResponse{
		Tutors: dtoList,
	})
	return nil
}
