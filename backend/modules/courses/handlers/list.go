package handlers

import (
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/courses/dtos"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type CourseListResponse struct {
	Items      []dtos.CourseListItemDTO `json:"items"`
	ItemsCount int64                    `json:"itemsCount"`
}

type CourseListRequest struct {
	common.SearchRequest
}

// @Summary List all available courses
// @Tags Courses
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body CourseListRequest true "Ability to filter results"
// @Success 200 {object} CourseListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses [get]
func List(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, _, _, searchParams := utils.GetRequestDataWithSearch[
		any,
		any,
	](c, "search")
	if err != nil {
		return err
	}

	// Check role
	if userData.Type != enums.UserTypeAdmin {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	m := models.Course{}

	var courses []models.Course
	query := initializers.DB.
		Model(&m)

	// Apply filters, sorting, pagination
	query, err = m.ApplyFilters(query, searchParams.ColumnFilters, m, map[string]interface{}{}, "")
	if err != nil {
		return err
	}
	query = m.ApplySorting(query, searchParams.Sorting, "year DESC, id DESC")
	totalCount := m.GetCount(query) // Gets count before pagination
	query = m.ApplyPagination(query, searchParams.Pagination)

	if err := query.Find(&courses).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch courses",
			Details: err.Error(),
		}
	}

	// Convert to DTOs
	dtoList := make([]dtos.CourseListItemDTO, len(courses))
	for i, q := range courses {
		dtoList[i] = dtos.CourseListItemDTO{}.From(&q)
	}

	c.JSON(200, CourseListResponse{
		Items:      dtoList,
		ItemsCount: totalCount,
	})

	return nil
}
