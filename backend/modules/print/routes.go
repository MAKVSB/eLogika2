package print

import (
	"elogika.vsb.cz/backend/modules/auth/wrappers"
	"elogika.vsb.cz/backend/modules/print/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("courses/:courseId/print/tests", wrappers.WithUserDataRole(handlers.PrintTest))
	rg.POST("courses/:courseId/print/questions", wrappers.WithUserDataRole(handlers.PrintQuestion))
	rg.POST("courses/:courseId/print/questions/:questionId", wrappers.WithUserDataRole(handlers.PrintQuestion))
}
