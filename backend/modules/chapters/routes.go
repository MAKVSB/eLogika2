package chapters

import (
	"elogika.vsb.cz/backend/modules/auth/wrappers"
	"elogika.vsb.cz/backend/modules/chapters/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("courses/:courseId/chapters/:parentId", wrappers.WithUserDataRole(handlers.ChapterInsert))
	rg.PUT("courses/:courseId/chapters/:chapterId", wrappers.WithUserDataRole(handlers.ChapterUpdate))
	rg.GET("courses/:courseId/chapters/:chapterId", wrappers.WithUserDataRole(handlers.ChapterGetByID))
	rg.GET("courses/:courseId/chapters", wrappers.WithUserDataRole(handlers.List))

	rg.PATCH("courses/:courseId/chapters/:chapterId/move/:direction", wrappers.WithUserDataRole(handlers.ChapterMove))
}
