package files

import (
	"elogika.vsb.cz/backend/modules/auth/wrappers"
	"elogika.vsb.cz/backend/modules/files/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("files", wrappers.WithUserData(handlers.FileUpload))
}

func RegisterRoutesUnauth(rg *gin.RouterGroup) {
	rg.GET("files/:fileId", handlers.FileServe)
}
