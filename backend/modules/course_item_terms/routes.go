package course_item_terms

import (
	"elogika.vsb.cz/backend/modules/auth/wrappers"
	"elogika.vsb.cz/backend/modules/course_item_terms/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("courses/:courseId/items/:courseItemId/terms", wrappers.WithUserDataRole(handlers.Insert))
	rg.PUT("courses/:courseId/items/:courseItemId/terms/:termId", wrappers.WithUserDataRole(handlers.Update))
	rg.DELETE("courses/:courseId/items/:courseItemId/terms/:termId", wrappers.WithUserDataRole(handlers.Delete))
	rg.GET("courses/:courseId/items/:courseItemId/terms", wrappers.WithUserDataRole(handlers.ListByItem))
	rg.GET("courses/:courseId/items/:courseItemId/terms/recursive", wrappers.WithUserDataRole(handlers.ListByItemRecursive))
	rg.GET("courses/:courseId/items/:courseItemId/terms/:termId", wrappers.WithUserDataRole(handlers.GetByID))

	rg.GET("courses/:courseId/items/:courseItemId/terms/:termId/students", wrappers.WithUserDataRole(handlers.ListJoinedStudents))
	rg.POST("courses/:courseId/items/:courseItemId/terms/:termId/students", wrappers.WithUserDataRole(handlers.UserJoin))
	rg.DELETE("courses/:courseId/items/:courseItemId/terms/:termId/students", wrappers.WithUserDataRole(handlers.UserLeave))

	rg.GET("courses/:courseId/terms", wrappers.WithUserDataRole(handlers.ListForStudent))
}
