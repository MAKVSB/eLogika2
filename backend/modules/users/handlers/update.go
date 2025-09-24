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

// @Description Request to create new user
type UserUpdateRequest struct {
	FirstName    string                   `json:"firstName" binding:"required"`    // First name
	FamilyName   string                   `json:"familyName" binding:"required"`   // (Family) Last name
	Username     string                   `json:"username" binding:"required"`     // Username
	Email        string                   `json:"email" binding:"required"`        // Email of the user
	Notification dtos.UserNotificationDTO `json:"notification" binding:"required"` // Notification setting
	Type         enums.UserTypeEnum       `json:"type" binding:"required"`         // System-wide user permissions
	Version      uint                     `json:"version" binding:"required"`      // Version signature to prevent concurrency problems
}

// @Description Newly created user
type UserUpdateResponse struct {
	Data dtos.UserDTO `json:"data"`
}

// @Summary Modify user
// @Tags Users
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param userId path int true "ID of the updated user"
// @Param body body UserUpdateRequest true "New data for user"
// @Success 200 {object} UserUpdateResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/users/{userId} [put]
func UserUpdate(c *gin.Context, userData authdtos.LoggedUserDTO) {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			UserID uint `uri:"userId" binding:"required"`
		},
		UserUpdateRequest,
	](c)
	if err != nil {
		c.AbortWithStatusJSON(err.Code, err)
		return
	}

	// TODO validate from here

	// check permissions
	if userData.Type != enums.UserTypeAdmin ||
		params.UserID != userData.ID {
		c.JSON(403, common.ErrorResponse{
			Message: "Not enough permissions",
		})
		return
	}

	transaction := initializers.DB.Begin()

	var user models.User
	if err := transaction.First(&user, params.UserID).Error; err != nil {
		transaction.Rollback()
		c.AbortWithStatusJSON(404, common.ErrorResponse{
			Message: "Failed to load user",
		})
		return
	}

	if user.Version != reqData.Version {
		transaction.Rollback()
		c.AbortWithStatusJSON(409, common.ErrorResponse{
			Message: "Version mismatched",
		})
		return
	}

	// Update only selected values
	user.Version = user.Version + 1
	if userData.Type == enums.UserTypeAdmin {
		user.Type = reqData.Type
	}
	if user.IdentityProvider == enums.IdentityProviderInternal {
		user.FirstName = reqData.FirstName
		user.FamilyName = reqData.FamilyName
		user.Username = reqData.Username
		user.Email = reqData.Email
	}
	user.Notification.Discord.Level.Messages = reqData.Notification.Discord.Level.Messages
	user.Notification.Discord.Level.Results = reqData.Notification.Discord.Level.Results
	user.Notification.Discord.Level.Terms = reqData.Notification.Discord.Level.Terms
	user.Notification.Discord.UserID = reqData.Notification.Discord.UserID
	user.Notification.Email.Level.Messages = reqData.Notification.Email.Level.Messages
	user.Notification.Email.Level.Results = reqData.Notification.Email.Level.Results
	user.Notification.Email.Level.Terms = reqData.Notification.Email.Level.Terms
	user.Notification.Push.Level.Messages = reqData.Notification.Push.Level.Messages
	user.Notification.Push.Level.Results = reqData.Notification.Push.Level.Results
	user.Notification.Push.Level.Terms = reqData.Notification.Push.Level.Terms

	if err := transaction.Save(&user).Error; err != nil {
		transaction.Rollback()
		c.AbortWithStatusJSON(500, common.ErrorResponse{
			Message: "Failed to update user",
		})
		return
	}

	if err := transaction.
		First(&user, user.ID).Error; err != nil {
		transaction.Rollback()
		c.AbortWithStatusJSON(500, common.ErrorResponse{
			Message: "Failed to fetch updated data",
		})
		return
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		c.AbortWithStatusJSON(500, common.ErrorResponse{
			Message: "Failed to commit changes",
		})
		return
	}

	c.JSON(200, UserUpdateResponse{
		Data: dtos.UserDTO{}.From(&user),
	})
}
