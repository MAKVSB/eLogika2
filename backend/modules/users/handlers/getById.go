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

// @Description User data
type UserGetByIdResponse struct {
	Data dtos.UserDTO `json:"data"`
}

// @Summary Get user by id
// @Tags Users
// @Security ApiKeyAuth
// @Produce  json
// @Param userId path int true "ID of the requested user"
// @Success 200 {object} UserGetByIdResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/users/{userId} [get]
func GetByID(c *gin.Context, userData authdtos.LoggedUserDTO) {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			UserID uint `uri:"userId" binding:"required"`
		},
		any,
	](c)
	if err != nil {
		c.AbortWithStatusJSON(err.Code, err)
		return
	}

	// TODO validate from here

	// check permissions
	if userData.Type != enums.UserTypeAdmin && params.UserID != userData.ID {
		c.JSON(403, common.ErrorResponse{
			Message: "Not enough permissions",
		})
		return
	}

	var user models.User
	if err := initializers.DB.
		Preload("UserCourses").
		Preload("UserCourses.Course").
		First(&user, params.UserID).Error; err != nil {
		c.AbortWithStatusJSON(404, common.ErrorResponse{
			Message: "Failed to fetch user",
		})
		return
	}

	c.JSON(200, UserGetByIdResponse{
		Data: dtos.UserDTO{}.From(&user),
	})
}
