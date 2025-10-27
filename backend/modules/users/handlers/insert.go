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

// @Description Request to insert new user
type UserInsertRequest struct {
	DegreeBefore string                   `json:"degreeBefore"`                    // Degree before name
	FirstName    string                   `json:"firstName" binding:"required"`    // First name
	FamilyName   string                   `json:"familyName" binding:"required"`   // (Family) Last name
	DegreeAfter  string                   `json:"degreeAfter"`                     // Degree after name
	Username     string                   `json:"username" binding:"required"`     // Username
	Email        string                   `json:"email" binding:"required"`        // Email of the user
	Notification dtos.UserNotificationDTO `json:"notification" binding:"required"` // Notification setting
	Type         enums.UserTypeEnum       `json:"type" binding:"required"`         // System-wide user permissions
}

// @Description Newly created user
type UserInsertResponse struct {
	Data dtos.UserDTO `json:"data"`
}

// @Summary Create new user
// @Tags Users
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body UserInsertRequest true "New data for user"
// @Success 200 {object} UserInsertResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/users/ [post]
func UserInsert(c *gin.Context, userData authdtos.LoggedUserDTO) {
	// Load request data
	err, _, reqData := utils.GetRequestData[
		any,
		UserInsertRequest,
	](c)
	if err != nil {
		c.AbortWithStatusJSON(err.Code, err)
		return
	}

	// TODO validate from here

	// check permissions
	if userData.Type != enums.UserTypeAdmin {
		c.JSON(403, common.ErrorResponse{
			Message: "Not enough permissions",
		})
		return
	}

	// TODO Generate random password and send it to email

	user := models.User{
		ID:           0,
		Version:      1,
		DegreeBefore: reqData.DegreeBefore,
		FirstName:    reqData.FirstName,
		FamilyName:   reqData.FamilyName,
		DegreeAfter:  reqData.DegreeAfter,
		Username:     reqData.Username,
		// Password:           "",
		Email:            reqData.Email,
		Type:             reqData.Type,
		IdentityProvider: enums.IdentityProviderInternal,
		Notification: models.UserNotification{
			Discord: models.NotificationDiscord{
				Level: models.NotificationLevel{
					Results:  reqData.Notification.Discord.Level.Results,
					Messages: reqData.Notification.Discord.Level.Messages,
					Terms:    reqData.Notification.Discord.Level.Terms,
				},
				UserID: reqData.Notification.Discord.UserID,
			},
			Email: models.NotificationEmail{
				Level: models.NotificationLevel{
					Results:  reqData.Notification.Email.Level.Results,
					Messages: reqData.Notification.Email.Level.Messages,
					Terms:    reqData.Notification.Email.Level.Terms,
				},
			},
			Push: models.NotificationPush{
				Level: models.NotificationLevel{
					Results:  reqData.Notification.Push.Level.Results,
					Messages: reqData.Notification.Push.Level.Messages,
					Terms:    reqData.Notification.Push.Level.Terms,
				},
				Token: "",
			},
		},
	}

	// TODO Generate random password and send it to email

	transaction := initializers.DB.Begin()

	if err := transaction.Save(&user).Error; err != nil {
		transaction.Rollback()
		c.AbortWithStatusJSON(500, common.ErrorResponse{
			Message: "Failed to insert user",
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

	c.JSON(200, UserInsertResponse{
		Data: dtos.UserDTO{}.From(&user),
	})
}
