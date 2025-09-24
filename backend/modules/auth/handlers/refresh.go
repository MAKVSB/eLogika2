package handlers

import (
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/auth/tokens"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"github.com/gin-gonic/gin"
)

// @Description Access token
type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}

// @Summary User refresh
// @Tags Auth
// @Produce  json
// @Success 200 {object} RefreshResponse "Successful operation"
// @Failure 401 {object} common.ErrorResponse "Unauthorised"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/auth/refresh [post]
func Refresh(c *gin.Context) {
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

	// Parse access token
	accessToken := tokens.AccessToken{}
	err = accessToken.Get(c, true)
	if err != nil {
		c.AbortWithStatusJSON(401, err)
		return
	}
	// I dont particulary caare if access token is expired here

	// Get updated user state and permissions
	var user models.User
	if err := initializers.DB.
		Preload("UserCourses.Course").
		Where("id = ?", refreshToken.UserID).First(&user).Error; err != nil {
		c.JSON(401, common.ErrorResponse{
			Message: "Invalid credentials",
		})
		return
	}

	courseService := services.NewCourseService(&repositories.CourseRepository{})
	courses, err := courseService.GetUserCourses(user.ID, &user.Type)
	if err != nil {
		c.JSON(403, common.ErrorResponse{
			Message: "Not enough permissions",
		})
		return
	}

	user.UserCourses = courses

	loggedUser := dtos.LoggedUserDTO{}.From(&user)

	// Revoke old and issue new access and refresh token
	newRefreshToken := &tokens.RefreshToken{}
	newAccessToken := &tokens.AccessToken{}
	newRefreshTokenStr, _ := newRefreshToken.New(loggedUser)
	newAccessTokenStr, _ := newAccessToken.New(loggedUser)
	refreshToken.Invalidate()
	accessToken.Invalidate()

	c.SetCookie(
		"refresh_token",
		newRefreshTokenStr,
		int(initializers.GlobalAppConfig.REFRESH_LENGTH.Seconds()),
		"/",
		"",
		false,
		true,
	)
	c.JSON(200, RefreshResponse{
		AccessToken: newAccessTokenStr,
	})
}
