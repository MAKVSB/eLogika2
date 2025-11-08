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
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Request to insert new course
type ClassUpdateRequest struct {
	Name          string                    `json:"name" binding:"required"`
	Room          string                    `json:"room" binding:"required"`
	Type          enums.ClassTypeEnum       `json:"type" binding:"required"`
	StudyForm     enums.StudyFormEnum       `json:"studyForm" binding:"required"`
	TimeFrom      string                    `json:"timeFrom" binding:"required"`
	TimeTo        string                    `json:"timeTo" binding:"required"`
	Day           enums.WeekDayEnum         `json:"day" binding:"required"`
	WeekParity    enums.WeekParityEnum      `json:"weekParity" binding:"required"`
	StudentLimit  uint                      `json:"studentLimit" binding:"required"`
	Version       uint                      `json:"version" binding:"required"` // Version signature to prevent concurrency problems
	ImportOptions models.ClassImportOptions `json:"importOptions" binding:"required"`
}

// @Description Newly created course
type ClassUpdateResponse struct {
	Data dtos.ClassDTO `json:"data"`
}

// @Summary Modify course
// @Tags Courses
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the updated course"
// @Param body body ClassUpdateRequest true "New data for course"
// @Success 200 {object} ClassUpdateResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId} [put]
func ClassUpdate(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
			ClassID  uint `uri:"classId" binding:"required"`
		},
		ClassUpdateRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	re := regexp.MustCompile(`^\d{2}:\d{2}$`)
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

	transaction := initializers.DB.Begin()

	classRepo := repositories.NewClassRepository()
	class, err := classRepo.GetClassByIDGarant(transaction, params.CourseID, params.ClassID, userData.ID, nil, true, &reqData.Version)
	if err != nil {
		transaction.Rollback()
		return err
	}

	// Update only selected values
	class.Version = reqData.Version + 1
	class.Name = reqData.Name
	class.Room = reqData.Room
	class.Type = reqData.Type
	class.StudyForm = reqData.StudyForm
	class.TimeFrom = reqData.TimeFrom
	class.TimeTo = reqData.TimeTo
	class.Day = reqData.Day
	class.WeekParity = reqData.WeekParity
	class.StudentLimit = reqData.StudentLimit
	class.ImportOptions = reqData.ImportOptions

	if err := transaction.Save(&class).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update class",
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	class, err = classRepo.GetClassByIDGarant(initializers.DB, params.CourseID, params.ClassID, userData.ID, nil, true, nil)
	if err != nil {
		return err
	}

	c.JSON(200, ClassUpdateResponse{
		Data: dtos.ClassDTO{}.From(class),
	})
	return nil
}

// TODO TEST
