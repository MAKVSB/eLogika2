package handlers

import (
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	authhelpers "elogika.vsb.cz/backend/modules/auth/helpers"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/users/helpers"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Request to create new user
type UserChangePassRequest struct {
	OldPassword *string `json:"oldPassword"` // First name
	NewPassword *string `json:"newPassword"` // (Family) Last name
	Generate    bool    `json:"generate"`    // Username
}

// @Description Newly created user
type UserChangePassResponse struct {
	Success bool `json:"success"`
}

// @Summary Modify or generate user password
// @Tags Users
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param userId path int true "ID of the updated user"
// @Param body body UserChangePassRequest true "New data for user"
// @Success 200 {object} UserChangePassResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/users/{userId}/password [put]
func ChangePass(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			UserID *uint `uri:"userId"`
		},
		UserChangePassRequest,
	](c)
	if err != nil {
		return err
	}

	if reqData.Generate {
		if userData.Type != enums.UserTypeAdmin {
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permissions",
			}
		}

		var user models.User
		if err := initializers.DB.First(&user, params.UserID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    404,
				Message: "Failed to load user",
				Details: "User not found",
			}
		}

		newPassword := helpers.GeneratePassword(&user)
		newPassHash, err := authhelpers.HashPassword(newPassword)
		if err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Details: "Failed to hash password",
			}
		}

		if err := initializers.DB.Model(&user).Where("id = ?", params.UserID).Update("password", newPassHash).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to save user data",
			}
		}
		// TODO send email
		// email := models.Email{
		// 	ID:       0,
		// 	ToEmail:  "",
		// 	Subject:  "",
		// 	Body:     "",
		// 	Status:   "",
		// 	Priority: 100,
		// }

		c.JSON(200, UserChangePassResponse{
			Success: true,
		})
		return nil
	} else {
		if params.UserID == nil {
			params.UserID = &userData.ID
		} else {
			if params.UserID != &userData.ID && userData.Type != enums.UserTypeAdmin {
				return &common.ErrorResponse{
					Code:    403,
					Message: "Not enough permissions",
				}
			}
		}

		var user models.User
		if err := initializers.DB.First(&user, userData.ID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    404,
				Message: "Failed to load user",
				Details: "User not found",
			}
		}

		if !authhelpers.CheckPasswordHash(*reqData.OldPassword, user.Password) {
			return &common.ErrorResponse{
				Code:    401,
				Message: "Invalid credentials",
				Details: "Old password does not match",
				FormErrors: common.ErrorObject{
					"oldPassword": "Password is incorrect",
				},
			}
		}

		passValid, passError := helpers.PasswordCheck(*reqData.NewPassword, &user)
		if !passValid {
			return &common.ErrorResponse{
				Code:    401,
				Message: "Invalid credentials",
				Details: "New password is not secure enough",
				FormErrors: common.ErrorObject{
					"newPassword":    passError,
					"newPasswordRep": passError,
				},
			}
		}

		newPassHash, err := authhelpers.HashPassword(*reqData.NewPassword)
		if err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Details: "Failed to hash password",
			}
		}

		if err := initializers.DB.Model(&user).Where("id = ?", user.ID).Update("password", newPassHash).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to save user data",
			}
		}

		c.JSON(200, UserChangePassResponse{
			Success: true,
		})
		return nil
	}
}
