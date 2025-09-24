package course_items

import (
	"elogika.vsb.cz/backend/modules/auth/wrappers"
	"elogika.vsb.cz/backend/modules/course_items/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("courses/:courseId/items", wrappers.WithUserDataRole(handlers.Insert))
	rg.PUT("courses/:courseId/items/:courseItemId", wrappers.WithUserDataRole(handlers.Update))
	rg.DELETE("courses/:courseId/items/:courseItemId", wrappers.WithUserDataRole(handlers.Delete))

	rg.GET("courses/:courseId/items", wrappers.WithUserDataRole(handlers.List))
	rg.GET("courses/:courseId/items/students", wrappers.WithUserDataRole(handlers.ListForStudent))
	rg.GET("courses/:courseId/items/:courseItemId", wrappers.WithUserDataRole(handlers.GetByID))

	rg.GET("courses/:courseId/items/:courseItemId/results", wrappers.WithUserDataRole(handlers.ListResults))
}
