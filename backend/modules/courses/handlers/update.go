package handlers

import (
	"encoding/json"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/courses/dtos"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Request to insert new course
type CourseUpdateRequest struct {
	Name          string                     `json:"name" binding:"required"`                          // Name of the course
	Content       json.RawMessage            `json:"content" binding:"required" ts_type:"JSONContent"` // Course text in json (Using TipTap editor format)
	Shortname     string                     `json:"shortname" binding:"required"`                     // Short name fort course
	Public        bool                       `json:"public"`                                           // Can any user join ?
	Year          uint                       `json:"year" binding:"required"`                          // Start year of academic year
	Semester      enums.SemesterEnum         `json:"semester" binding:"required"`                      // Semester of the above year
	PointsMin     float64                    `json:"pointsMin" binding:"required"`                     // Minimum required points to pass
	PointsMax     float64                    `json:"pointsMax" binding:"required"`                     // Maximum points
	ImportOptions models.CourseImportOptions `json:"importOptions" binding:"required"`
	Version       uint                       `json:"version" binding:"required"` // Version signature to prevent concurrency problems
}

// @Description Newly created course
type CourseUpdateResponse struct {
	Data dtos.CourseDTO `json:"data"`
}

// @Summary Modify course
// @Tags Courses
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the updated course"
// @Param body body CourseUpdateRequest true "New data for course"
// @Success 200 {object} CourseUpdateResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId} [put]
func Update(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		CourseUpdateRequest,
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
	if userRole != enums.CourseUserRoleAdmin && userRole != enums.CourseUserRoleGarant {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	courseService := services.CourseService{}
	course, err := courseService.GetCourseByID(initializers.DB, params.CourseID, userData.ID, userRole, nil, true, &reqData.Version)
	if err != nil {
		return err
	}

	transaction := initializers.DB.Begin()

	// Update only selected values
	course.Version = reqData.Version + 1
	course.Name = reqData.Name
	course.Content = reqData.Content
	course.Shortname = reqData.Shortname
	course.Public = reqData.Public
	course.Year = reqData.Year
	course.PointsMin = reqData.PointsMin
	course.PointsMax = reqData.PointsMax
	course.Semester = reqData.Semester

	course.ImportOptions = reqData.ImportOptions

	if err := transaction.Save(&course).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update course",
		}
	}

	// Sync content courseFiles
	courseRepo := repositories.NewCourseRepository()
	err = courseRepo.SyncFiles(transaction, course.Content, course)
	if err != nil {
		transaction.Rollback()
		return err
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	if err := initializers.DB.
		First(&course, course.ID).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch updated data",
		}
	}

	c.JSON(200, CourseUpdateResponse{
		Data: dtos.CourseDTO{}.From(course),
	})

	return nil
}
