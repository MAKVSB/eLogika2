package handlers

import (
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/auth/helpers"
	"elogika.vsb.cz/backend/modules/auth/tokens"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`    // User email
	Password string `json:"password" binding:"required"` // User password
}

// @Description Access token
type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

// @Summary User login
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param body body LoginRequest true "Login request data"
// @Success 200 {object} LoginResponse "Successful operation"
// @Failure 401 {object} common.ErrorResponse "Unauthorised"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/auth/login [post]
func Login(c *gin.Context) {
	// Load request data
	err, _, reqData := utils.GetRequestData[
		any,
		LoginRequest,
	](c)
	if err != nil {
		c.AbortWithStatusJSON(err.Code, err)
		return
	}

	// TODO validate from here

	var user models.User
	if err := initializers.DB.
		Where("email = ?", reqData.Email).First(&user).Error; err != nil {
		c.JSON(401, common.ErrorResponse{
			Message: "Invalid credentials",
		})
		return
	}

	if !helpers.CheckPasswordHash(reqData.Password, user.Password) {
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

	accessToken := tokens.AccessToken{}
	refreshToken := tokens.RefreshToken{}
	accessTokenStr, _ := accessToken.New(loggedUser)
	refreshTokenStr, _ := refreshToken.New(loggedUser)

	c.SetCookie(
		"refresh_token",
		refreshTokenStr,
		int(initializers.GlobalAppConfig.REFRESH_LENGTH.Seconds()),
		"/",
		"",
		false,
		true,
	)
	c.JSON(200, LoginResponse{
		AccessToken: accessTokenStr,
	})
}
