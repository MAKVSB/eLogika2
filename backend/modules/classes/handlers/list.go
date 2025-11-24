package handlers

import (
	"elogika.vsb.cz/backend/initializers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/classes/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
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
// @Router /api/v2/courses/{courseId}/classes [get]
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

	classService := services.NewClassService(repositories.NewClassRepository())
	classes, classCount, err := classService.ListClasses(initializers.DB, params.CourseID, userData.ID, userRole, nil, true, searchParams)
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
