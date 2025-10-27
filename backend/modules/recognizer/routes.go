package recognizer

import (
	"elogika.vsb.cz/backend/modules/auth/wrappers"
	"elogika.vsb.cz/backend/modules/recognizer/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("recognizer/test/:identifier", wrappers.WithUserDataRole(handlers.RecognizerTestGet))
	rg.POST("recognizer/test", wrappers.WithUserDataRole(handlers.RecognizerTestSave))

}
