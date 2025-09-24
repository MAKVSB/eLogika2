package handlers

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Description Request to upload a file
type FileUploadRequest struct {
	ReferencingEntityType string `json:"referencingEntityType"` // Type of the resource the file is connected to
	ReferencingEntity     uint   `json:"referencingEntity"`     // ID of the resource the file is connected to
}

// @Description Reference to uploaded file
type FileUploadResponse struct {
	ID           uint   `json:"id"`
	Filename     string `json:"filename"`
	OriginalName string `json:"originalName"`
}

// @Summary Upload new file
// @Tags Files
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body FileUploadRequest true "New data for chapter"
// @Success 200 {object} FileUploadResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/files [post]
func FileUpload(c *gin.Context, userData authdtos.LoggedUserDTO) {
	// Load request data
	err, _, _ := utils.GetRequestData[
		any,
		any,
	](c)
	if err != nil {
		c.AbortWithStatusJSON(err.Code, err)
		return
	}

	// TODO validate from here

	// TODO Figure out if it is good idea to allow everyone upload everything without any checks
	// TODO For example. Tutor can try to upload image before the question even exists.
	// TODO It can probably happen that the file will be uploaded, but not linked to anything ???
	// TODO add check of permissions if it is at least tutor when trying to link to question
	// TODO only student is allowed to upload homeworks and other assignments

	// Handle file upload itself
	file, header, err2 := c.Request.FormFile("file")
	if err2 != nil {
		c.AbortWithStatusJSON(400, common.ErrorResponse{
			Message: "Invalid file",
		})
		return
	}
	defer file.Close()

	// Check extensions
	// TODO Introduce a list of allowed extensions
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" || len(ext) > 10 {
		c.AbortWithStatusJSON(400, common.ErrorResponse{
			Message: "Invalid file extension",
		})
		return
	}

	// Generate a random UUID filename
	newFileName := uuid.New().String() + ext
	for i := 0; i > 0; i++ {
		candidate := uuid.New().String() + ext

		var count int64
		if err := initializers.DB.Model(&models.File{}).
			Where("stored_name = ?", candidate).
			Count(&count).Error; err != nil {
			c.AbortWithStatusJSON(500, common.ErrorResponse{
				Message: "DB error during UUID check",
			})
			return
		}

		if count == 0 {
			newFileName = candidate
			break
		}
	}

	// Create file
	out, err2 := os.Create(filepath.Join(initializers.GlobalAppConfig.UPLOADS_DESTINATION, newFileName))
	if err2 != nil {
		c.AbortWithStatusJSON(500, common.ErrorResponse{
			Message: "Could not save file",
		})
		return
	}
	defer out.Close()

	// Copy data
	size, err2 := io.Copy(out, file)
	if err2 != nil {
		c.AbortWithStatusJSON(500, common.ErrorResponse{
			Message: "Failed to save file",
		})
		return
	}

	// Save to database for later linking
	uploaded := models.File{
		UserID:       userData.ID,
		OriginalName: header.Filename,
		StoredName:   newFileName,
		MIMEType:     header.Header.Get("Content-Type"),
		SizeBytes:    size,
		UploadedAt:   time.Now(),
	}

	if err := initializers.DB.Create(&uploaded).Error; err != nil {
		c.AbortWithStatusJSON(500, common.ErrorResponse{
			Message: "DB save failed",
		})
		return
	}

	c.JSON(200, FileUploadResponse{
		ID:           uploaded.ID,
		Filename:     uploaded.StoredName,
		OriginalName: uploaded.OriginalName,
	})
}
