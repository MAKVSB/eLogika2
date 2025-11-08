package handlers

import (
	"time"

	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/auth/tokens"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Request to generate new api token
type TokenCreateRequest struct {
	Name      string    `json:"name" binding:"required"`
	ExpiresAt time.Time `json:"expiresAt" binding:"required"`
}

// @Description Newly generated token
type TokenCreateResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// @Summary Create new api token
// @Tags Users
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body TokenCreateRequest true "New token"
// @Success 200 {object} TokenCreateResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/users/tokens [post]
func TokenCreate(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, _, reqData := utils.GetRequestData[
		any,
		TokenCreateRequest,
	](c)
	if err != nil {
		return err
	}

	// Creates token
	token := tokens.ApiToken{}
	apiTokenString, err2 := token.New(userData, reqData.Name, reqData.ExpiresAt)
	if err2 != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to generate token",
			Details: err2,
		}
	}

	// Response
	c.JSON(200, TokenCreateResponse{
		Name:  reqData.Name,
		Value: apiTokenString,
	})
	return nil
}
