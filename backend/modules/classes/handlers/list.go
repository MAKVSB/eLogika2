package handlers

import (
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

type ClassListResponse struct {
	Items      []dtos.ClassListItemDTO `json:"items"`
	ItemsCount int64                   `json:"itemsCount"`
}

type ClassListRequest struct {
	common.SearchRequest
}

// @Summary List classes thaat belongs to course
// @Tags Classes
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body ClassListRequest true "Ability to filter results"
// @Success 200 {object} ClassListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses [get]
func List(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _, searchParams := utils.GetRequestDataWithSearch[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		any,
	](c, "search")
	if err != nil {
		return err
	}

	// TODO validate from here

	classRepo := repositories.NewClassRepository()
	var classes []*models.Class
	var classCount int64
	// If not admin, garant or tutor
	if userRole == enums.CourseUserRoleAdmin {
		classes, classCount, err = classRepo.ListClassesAdmin(initializers.DB, params.CourseID, userData.ID, true, searchParams)
	} else if userRole == enums.CourseUserRoleGarant {
		classes, classCount, err = classRepo.ListClassesGarant(initializers.DB, params.CourseID, userData.ID, true, searchParams)
	} else if userRole == enums.CourseUserRoleTutor {
		classes, classCount, err = classRepo.ListClassesTutor(initializers.DB, params.CourseID, userData.ID, true, searchParams)
	} else {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
	if err != nil {
		return err
	}

	// Convert to DTOs
	dtoList := make([]dtos.ClassListItemDTO, len(classes))
	for i, c := range classes {
		dtoList[i] = dtos.ClassListItemDTO{}.From(c)
	}

	c.JSON(200, ClassListResponse{
		Items:      dtoList,
		ItemsCount: classCount,
	})
	return nil
}
