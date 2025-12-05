package questions

import (
	"elogika.vsb.cz/backend/modules/auth/wrappers"
	"elogika.vsb.cz/backend/modules/questions/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("courses/:courseId/questions", wrappers.WithUserDataRole(handlers.List))

	rg.POST("courses/:courseId/questions", wrappers.WithUserDataRole(handlers.QuestionInsert))
	rg.PUT("courses/:courseId/questions/:questionId", wrappers.WithUserDataRole(handlers.QuestionUpdate))
	rg.GET("courses/:courseId/questions/:questionId", wrappers.WithUserDataRole(handlers.QuestionGetByID))
	rg.DELETE("courses/:courseId/questions/:questionId", wrappers.WithUserDataRole(handlers.QuestionDelete))

	// Question checks
	rg.POST("courses/:courseId/questions/:questionId/check", wrappers.WithUserDataRole(handlers.Check))
	rg.DELETE("courses/:courseId/questions/:questionId/check", wrappers.WithUserDataRole(handlers.Uncheck))
	rg.PATCH("courses/:courseId/questions/:questionId/toggleActive", wrappers.WithUserDataRole(handlers.QuestionToggleActive))

	// Version change
	rg.PATCH("courses/:courseId/questions/:questionId/selectversion", wrappers.WithUserDataRole(handlers.SelectVersion))
}
