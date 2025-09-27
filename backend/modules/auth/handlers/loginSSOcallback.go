package handlers

import (
	"fmt"
	"net/url"
	"strings"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/auth/tokens"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/cas.v2"
)

// @Description Login request
type SSOLoginCallbackRequest struct {
	SSOProvider string `json:"provider"`
	Ticket      string `json:"ticket"`
}

// @Description Access token
type SSOLoginCallbackResponse struct {
	AccessToken string `json:"accessToken"`
}

// @Summary User login
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param body body SSOLoginCallbackRequest true "Login request data"
// @Success 200 {object} SSOLoginCallbackResponse "Successful operation"
// @Failure 401 {object} common.ErrorResponse "Unauthorised"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/auth/login [post]
func SSOLoginCallback(c *gin.Context) {
	// Load request data
	err, _, reqData := utils.GetRequestData[
		any,
		SSOLoginCallbackRequest,
	](c)
	if err != nil {
		c.AbortWithStatusJSON(err.Code, err)
		return
	}

	if reqData.SSOProvider == "VSBCAS" {
		// request from VSBCAS
		casUrl, err := url.Parse("https://www.sso.vsb.cz/")
		if err != nil {
			panic("Failed to parse VSB-SSO url")
		}

		serviceUrl, err := url.Parse("https://elogika.vsb.cz/api/v2/auth/login/sso")
		if err != nil {
			panic("Failed to parse VSB-SSO url")
		}

		casRestClient := cas.NewRestClient(&cas.RestOptions{
			CasURL:     casUrl,
			ServiceURL: serviceUrl,
		})

		test1, err := casRestClient.ValidateServiceTicket(cas.ServiceTicket(reqData.Ticket))
		if err != nil {
			fmt.Println(err.Error())
			utils.DebugPrintJSON(err)
			errr := &common.ErrorResponse{
				Code:    500,
				Message: "Network error during validation",
			}
			c.AbortWithStatusJSON(errr.Code, errr)
			return
		}

		username := test1.User
		personId := test1.Attributes.Get("personId")
		if username == "" || personId == "" {
			errr := &common.ErrorResponse{
				Code:    500,
				Message: "Token is missing important data",
			}
			c.AbortWithStatusJSON(errr.Code, errr)
			return
		}

		utils.DebugPrintJSON(username)
		utils.DebugPrintJSON(personId)

		var user models.User
		if err := initializers.DB.
			Where("username = ? AND identity_provider_id = ?", strings.ToUpper(username), personId).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(401, common.ErrorResponse{
				Message: "Invalid credentials",
			})
			return
		}

		courseService := services.NewCourseService(&repositories.CourseRepository{})
		courses, err2 := courseService.GetUserCourses(user.ID, &user.Type)
		if err2 != nil {
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
	} else {
		errr := &common.ErrorResponse{
			Code:    500,
			Message: "Implemented providers are: \"VSBCAS\"",
		}
		c.AbortWithStatusJSON(errr.Code, errr)
		return
	}
}
