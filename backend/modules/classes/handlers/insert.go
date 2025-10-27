package handlers

import (
	"regexp"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/classes/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Request to insert new class
type ClassInsertRequest struct {
	Name          string                    `json:"name" binding:"required"`
	Room          string                    `json:"room" binding:"required"`
	Type          enums.ClassTypeEnum       `json:"type" binding:"required"`
	StudyForm     enums.StudyFormEnum       `json:"studyForm" binding:"required"`
	TimeFrom      string                    `json:"timeFrom" binding:"required"`
	TimeTo        string                    `json:"timeTo" binding:"required"`
	Day           enums.WeekDayEnum         `json:"day" binding:"required"`
	WeekParity    enums.WeekParityEnum      `json:"weekParity" binding:"required"`
	StudentLimit  uint                      `json:"studentLimit" binding:"required"`
	ImportOptions models.ClassImportOptions `json:"importOptions" binding:"required"`
}

// @Description Newly created course
type ClassInsertResponse struct {
	Data dtos.ClassDTO `json:"data"`
}

// @Summary Create new class
// @Tags Classes
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body ClassInsertRequest true "New data for class"
// @Success 200 {object} ClassInsertResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/ [post]
func ClassInsert(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		ClassInsertRequest,
	](c)
	if err != nil {
		return err
	}
	// TODO validate from here

	// Validate other inputs
	re := regexp.MustCompile(`^\d{1,2}:\d{2}$`)
	if !re.MatchString(reqData.TimeFrom) {
		return &common.ErrorResponse{
			Code:    422,
			Message: "Incorrect time timeFrom format",
		}
	}
	if !re.MatchString(reqData.TimeTo) {
		return &common.ErrorResponse{
			Code:    422,
			Message: "Incorrect time timeTo format",
		}
	}

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}
	// If not admin or garant
	if userRole != enums.CourseUserRoleAdmin && userRole != enums.CourseUserRoleGarant {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	class := models.Class{
		ID:            0,
		Version:       1,
		CourseID:      params.CourseID,
		Name:          reqData.Name,
		Room:          reqData.Room,
		Type:          reqData.Type,
		StudyForm:     reqData.StudyForm,
		TimeFrom:      reqData.TimeFrom,
		TimeTo:        reqData.TimeTo,
		Day:           reqData.Day,
		WeekParity:    reqData.WeekParity,
		StudentLimit:  reqData.StudentLimit,
		ImportOptions: reqData.ImportOptions,
	}

	transaction := initializers.DB.Begin()

	if err := transaction.Save(&class).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to insert class",
		}
	}

	if err := transaction.
		First(&class, class.ID).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch updated data",
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	c.JSON(200, ClassInsertResponse{
		Data: dtos.ClassDTO{}.From(&class),
	})
	return nil
}

// TODO TEST
