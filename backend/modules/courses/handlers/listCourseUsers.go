package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/courses/dtos"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Updated list of course users
type ListCourseUsersResponse struct {
	Items      []dtos.CourseUserDTO `json:"items"`
	ItemsCount uint                 `json:"itemsCount"`
}

// @Summary List tutors of class
// @Tags Courses
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} ListCourseUsersResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/users [get]
func ListCourseUsers(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
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

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}
	// If not admin, or garant
	if userRole != enums.CourseUserRoleTutor && userRole != enums.CourseUserRoleAdmin && userRole != enums.CourseUserRoleGarant {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	innerQuery := initializers.DB.Model(models.User{})
	innerQuery, err = models.User{}.ApplyFilters(innerQuery, searchParams.ColumnFilters, models.User{}, map[string]interface{}{}, "")
	if err != nil {
		return err
	}
	innerQuery = models.User{}.ApplySorting(innerQuery, searchParams.Sorting, "")

	query := initializers.DB.Model(models.CourseUser{}).
		Where("course_id = ?", params.CourseID).
		InnerJoins("User", innerQuery)

	query, err = models.CourseUser{}.ApplyFilters(query, searchParams.ColumnFilters, models.CourseUser{}, map[string]interface{}{}, "")
	if err != nil {
		return err
	}
	query = models.CourseUser{}.ApplySorting(query, searchParams.Sorting, "")
	totalCount := models.CourseUser{}.GetCount(query) // Gets count before pagination
	query = models.CourseUser{}.ApplyPagination(query, searchParams.Pagination)

	var courseUsers []models.CourseUser
	if err := query.
		Find(&courseUsers).Error; err != nil {
		return &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch users 2",
			Details: err.Error(),
		}
	}

	// Convert to DTOs
	dtoList := make([]dtos.CourseUserDTO, len(courseUsers))
	for i, c := range courseUsers {
		dtoList[i] = dtos.CourseUserDTO{}.From(&c)
	}

	utils.DebugPrintJSON(searchParams)

	c.JSON(200, ListCourseUsersResponse{
		Items:      dtoList,
		ItemsCount: uint(totalCount),
	})

	return nil
}
