package handlers

import (
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	authEnums "elogika.vsb.cz/backend/modules/auth/enums"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/users/dtos"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Description User data
type UserGetSelfResponse struct {
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
// @Router /api/v2/users/self [get]
func GetSelf(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, _, _ := utils.GetRequestData[
		any,
		any,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	var user models.User
	if err := initializers.DB.
		Preload("ApiTokens", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("token_type = ?", authEnums.JWTTokenTypeApi).
				Where("expires_at > ?", time.Now()).
				Where("revoked_at IS NULL")
		}).
		First(&user, userData.ID).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch user",
		}
	}

	c.JSON(200, UserGetByIdResponse{
		Data: dtos.UserDTO{}.From(&user),
	})

	return nil
}
