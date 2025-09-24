package categories

import (
	"elogika.vsb.cz/backend/modules/auth/wrappers"
	"elogika.vsb.cz/backend/modules/categories/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("courses/:courseId/chapters/:chapterId/categories", wrappers.WithUserDataRole(handlers.List))

	rg.POST("courses/:courseId/categories", wrappers.WithUserDataRole(handlers.Insert))
	rg.PUT("courses/:courseId/categories/:categoryId", wrappers.WithUserDataRole(handlers.Update))
	rg.GET("courses/:courseId/categories/:categoryId", wrappers.WithUserDataRole(handlers.GetByID))
	rg.GET("courses/:courseId/categories", wrappers.WithUserDataRole(handlers.List))
}
