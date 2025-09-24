package handlers

import (
	"elogika.vsb.cz/backend/modules/auth/tokens"
	"elogika.vsb.cz/backend/modules/common"
	"github.com/gin-gonic/gin"
)

// @Summary User logout from all devices and sessions
// @Tags Auth
// @Produce  json
// @Success 200 {object} LogoutResponse "Successful operation"
// @Failure 401 {object} common.ErrorResponse "Unauthorised"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/auth/logout/all [post]
func LogoutAll(c *gin.Context, accesstoken tokens.Token) {
	// Parse refresh token
	refreshToken := tokens.RefreshToken{}
	err := refreshToken.Get(c, false)
	if err != nil {
		c.AbortWithStatusJSON(401, err)
		return
	}
	if refreshToken.IsRevoked() {
		c.AbortWithStatusJSON(401, common.ErrorResponse{
			Message: "Refresh token revoked",
		})
		return
	}

	// Revoke tokens
	refreshToken.InvalidateByUser()
	accesstoken.InvalidateByUser()

	c.JSON(200, LogoutResponse{
		Success: true,
	})
}
