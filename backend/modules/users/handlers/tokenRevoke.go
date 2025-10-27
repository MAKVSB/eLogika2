package handlers

import (
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"

	authenums "elogika.vsb.cz/backend/modules/auth/enums"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type TokenRevokeResponse struct {
	Success bool `json:"success"`
}

// @Summary Revoke token
// @Tags Users
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} TokenRevokeResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/users/self/tokens/:tokenId [delete]
func TokenRevoke(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			TokenID *string `uri:"tokenId"`
		},
		any,
	](c)
	if err != nil {
		return err
	}

	// Revokes token
	if err := initializers.DB.
		Model(&models.AuthToken{}).
		Where("token_id", params.TokenID).
		Where("token_type", authenums.JWTTokenTypeApi).
		Where("user_id", userData.ID).
		Update("revoked_at", time.Now()).
		Update("revoked_for", authenums.RevokedForToken).Error; err != nil {
		return &common.ErrorResponse{
			Code:    404,
			Message: "Failed to revoke token",
			Details: err.Error(),
		}
	}

	// Response
	c.JSON(200, TokenRevokeResponse{
		Success: true,
	})
	return nil
}
