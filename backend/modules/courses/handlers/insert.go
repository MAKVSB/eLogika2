package handlers

import (
	"encoding/json"

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
type CourseInsertRequest struct {
	Name          string                     `json:"name" binding:"required"`                          // Name of the course
	Content       json.RawMessage            `json:"content" binding:"required" ts_type:"JSONContent"` // Course text in json (Using TipTap editor format)
	Shortname     string                     `json:"shortname" binding:"required"`                     // Short name fort course
	Public        bool                       `json:"public"`                                           // Can any user join ?
	Year          uint                       `json:"year" binding:"required"`                          // Start year of academic year
	Semester      enums.SemesterEnum         `json:"semester" binding:"required"`                      // Semester of the above year
	ImportOptions models.CourseImportOptions `json:"importOptions" binding:"required"`
}

// @Description Newly created course
type CourseInsertResponse struct {
	Data dtos.CourseDTO `json:"data"`
}

// @Summary Create new course
// @Tags Courses
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body CourseInsertRequest true "New data for course"
// @Success 200 {object} CourseInsertResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/ [post]
func CourseInsert(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, _, reqData := utils.GetRequestData[
		any,
		CourseInsertRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role
	if userData.Type != enums.UserTypeAdmin {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	course := &models.Course{
		ID:            0,
		Version:       1,
		Name:          reqData.Name,
		Content:       reqData.Content,
		Shortname:     reqData.Shortname,
		Public:        reqData.Public,
		Year:          reqData.Year,
		Semester:      reqData.Semester,
		ImportOptions: reqData.ImportOptions,
	}

	transaction := initializers.DB.Begin()

	if err := transaction.Save(&course).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to insert course",
		}
	}

	// Sync content courseFiles
	courseRepo := repositories.NewCourseRepository()
	err = courseRepo.SyncFiles(transaction, course.Content, course)
	if err != nil {
		transaction.Rollback()
		return err
	}

	chapter := &models.Chapter{
		ID:       0,
		Version:  1,
		CourseID: course.ID,
		Name:     reqData.Name,
		Content:  reqData.Content,
		Visible:  false,
		Order:    1,
	}

	if err := transaction.Save(&chapter).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to insert chapter",
		}
	}

	// Link back the chapter ID
	course.ChapterID = &chapter.ID
	if err := transaction.Save(&course).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to insert course",
		}
	}

	// Sync content chapterFiles
	chapterRepo := repositories.NewChapterRepository()
	err = chapterRepo.SyncFiles(transaction, chapter.Content, chapter)
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

	courseService := services.NewCourseService(repositories.NewCourseRepository())
	course, err = courseService.GetCourseByID(initializers.DB, course.ID, userData.ID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	c.JSON(200, CourseInsertResponse{
		Data: dtos.CourseDTO{}.From(course),
	})
	return nil
}
