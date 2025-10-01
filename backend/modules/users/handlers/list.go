package handlers

import (
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/users/dtos"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Items      []dtos.UserListItemDTO `json:"items"`
	ItemsCount int64                  `json:"itemsCount"`
}

type UserListRequest struct {
	common.SearchRequest
}

// @Summary List all available users
// @Tags Users
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body UserListRequest true "Ability to filter results"
// @Success 200 {object} UserListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/users [get]
func List(c *gin.Context, userData authdtos.LoggedUserDTO) {
	// Load request data
	err, _, _, searchParams := utils.GetRequestDataWithSearch[
		any,
		any,
	](c, "search")
	if err != nil {
		c.AbortWithStatusJSON(err.Code, err)
		return
	}

	// TODO validate from here

	m := models.User{}

	// check permissions
	if userData.Type != enums.UserTypeAdmin {
		c.JSON(403, common.ErrorResponse{
			Message: "Not enough permissions",
		})
		return
	}

	var users []models.User
	query := initializers.DB.
		Model(&m)

	// Apply filters, sorting, pagination
	query, err = m.ApplyFilters(query, searchParams.ColumnFilters, m, map[string]interface{}{}, "")
	if err != nil {
		c.AbortWithStatusJSON(400, err)
		return
	}
	query = m.ApplySorting(query, searchParams.Sorting)
	totalCount := m.GetCount(query) // Gets count before pagination
	query = m.ApplyPagination(query, searchParams.Pagination)

	if err := query.Find(&users).Error; err != nil {
		c.AbortWithStatusJSON(500, common.ErrorResponse{
			Message: "Failed to fetch users",
			Details: err.Error(),
		})
		return
	}

	// Convert to DTOs
	dtoList := make([]dtos.UserListItemDTO, len(users))
	for i, q := range users {
		dtoList[i] = dtos.UserListItemDTO{}.From(&q)
	}

	c.JSON(200, UserListResponse{
		Items:      dtoList,
		ItemsCount: totalCount,
	})
}
