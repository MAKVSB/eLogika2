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

// @Description Updated list of tutors
type ListStudentResponse struct {
	Items      []dtos.ClassUserDTO `json:"items"`
	ItemsCount int64               `json:"itemsCount"`
}

// @Summary List tutors of class
// @Tags Classes
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} ListStudentResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/ [post] // TODO
func ListStudents(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _, searchParams := utils.GetRequestDataWithSearch[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
			ClassID  uint `uri:"classId" binding:"required"`
		},
		any,
	](c, "search")
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		return err
	}

	classRepo := repositories.NewClassRepository()
	var class *models.Class
	// If not admin, garant or tutor
	if userRole == enums.CourseUserRoleAdmin {
		class, err = classRepo.GetClassByIDAdmin(initializers.DB, params.CourseID, params.ClassID, userData.ID, false, nil)
	} else if userRole == enums.CourseUserRoleGarant {
		class, err = classRepo.GetClassByIDGarant(initializers.DB, params.CourseID, params.ClassID, userData.ID, false, nil)
	} else if userRole == enums.CourseUserRoleTutor {
		class, err = classRepo.GetClassByIDTutor(initializers.DB, params.CourseID, params.ClassID, userData.ID, false, nil)
	} else {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
	if err != nil {
		return err
	}
	if class == nil {
		panic("class must be found")
	}

	// TODO add pagination

	students, studentCount, err := classRepo.ListClassStudents(initializers.DB, params.CourseID, params.ClassID, userData.ID, nil, false, searchParams)
	if err != nil {
		return err
	}

	// Convert to DTOs
	dtoList := make([]dtos.ClassUserDTO, len(students))
	for i, c := range students {
		dtoList[i] = dtos.ClassUserDTO{}.From(c)
	}

	c.JSON(200, ListStudentResponse{
		Items:      dtoList,
		ItemsCount: studentCount,
	})
	return nil
}
