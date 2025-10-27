package handlers

import (
	"path/filepath"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Gets a file
// @Tags Files
// @Security ApiKeyAuth
// @Param fileId path int true "ID of the requested file"
// @Success 200 {object} FileUploadResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/files/{fileId} [get]
func FileServe(c *gin.Context) {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			FileID string `uri:"fileId" binding:"required"`
		},
		any,
	](c)
	if err != nil {
		c.AbortWithStatusJSON(err.Code, err)
		return
	}

	// Check if files really exists within uploaded files
	var fileData *models.File
	if err := initializers.DB.
		Where("stored_name = ?", params.FileID).
		Find(&fileData).Error; err != nil {
		errCode := &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch file",
		}
		c.AbortWithStatusJSON(errCode.Code, errCode)
		return
	}

	// TODO validate from here

	// TODO Check route signature

	filepath := filepath.Join(initializers.GlobalAppConfig.UPLOADS_DESTINATION, params.FileID)
	c.File(filepath)
}
