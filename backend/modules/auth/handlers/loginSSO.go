package handlers

import (
	"net/url"

	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Login request
type SSOLoginRequest struct {
	SSOProvider string `json:"provider"`
}

// @Description Access token
type SSOLoginResponse struct {
	RedirectUrl string `json:"redirectUrl"`
}

// @Summary User login
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param body body SSOLoginRequest true "Login request data"
// @Success 200 {object} SSOLoginResponse "Successful operation"
// @Failure 401 {object} common.ErrorResponse "Unauthorised"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/auth/login [post]
func SSOLogin(c *gin.Context) {
	// Load request data
	err, _, reqData := utils.GetRequestData[
		any,
		SSOLoginRequest,
	](c)
	if err != nil {
		c.AbortWithStatusJSON(err.Code, err)
		return
	}

	if reqData.SSOProvider == "VSBCAS" {
		params := url.Values{}
		params.Add("service", "https://elogika.vsb.cz/new/login/callback?provider=VSBCAS")

		c.JSON(200, SSOLoginResponse{
			RedirectUrl: "https://www.sso.vsb.cz/login?" + params.Encode(),
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
