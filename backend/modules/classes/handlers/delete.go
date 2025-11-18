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

// @Description Course data
type ClassDeleteResponse struct {
	Success bool `json:"success"`
}

// @Summary Get class by id
// @Tags Classes
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the requested course"
// @Success 200 {object} ClassGetByIdResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/classes/{classId} [delete]
func ClassDelete(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
			ClassID  uint `uri:"classId" binding:"required"`
		},
		any,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}

	transaction := initializers.DB.Begin()

	classService := services.NewClassService(repositories.NewClassRepository())
	class, err := classService.GetClassByID(transaction, params.CourseID, params.ClassID, userData.ID, userRole, nil, false, nil)
	if err != nil {
		transaction.Rollback()
		return err
	}

	// Check if class has students
	// Get term count
	var studentCount int64
	if err := transaction.
		Model(&models.ClassStudent{}).
		Where("class_id = ?", params.ClassID).
		Count(&studentCount).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to get students count",
		}
	}
	if studentCount != 0 {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    409,
			Message: "Failed to delete class",
			Details: "Class has students assigned",
		}
	}

	// Unlink class tutors
	var tutors []*models.ClassTutor
	if err := transaction.
		Where("class_id = ?", params.ClassID).
		Delete(&tutors).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to get unlink tutors",
		}
	}

	if err := transaction.Delete(&class).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to get delete class",
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to get commit changes",
		}
	}

	c.JSON(200, ClassGetByIdResponse{
		Data: dtos.ClassDTO{}.From(class),
	})
	return nil
}
