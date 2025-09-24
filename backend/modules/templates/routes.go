package templates

import (
	"elogika.vsb.cz/backend/modules/auth/wrappers"
	"elogika.vsb.cz/backend/modules/templates/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("courses/:courseId/templates", wrappers.WithUserDataRole(handlers.List))
	rg.POST("courses/:courseId/templates", wrappers.WithUserDataRole(handlers.Insert))

	rg.GET("courses/:courseId/templates/creator", wrappers.WithUserDataRole(handlers.Creator))
	rg.GET("courses/:courseId/templates/:templateId", wrappers.WithUserDataRole(handlers.TemplateGetByID))
	rg.PUT("courses/:courseId/templates/:templateId", wrappers.WithUserDataRole(handlers.Update))
	rg.DELETE("courses/:courseId/templates/:templateId", wrappers.WithUserDataRole(handlers.TemplateDelete))
}
