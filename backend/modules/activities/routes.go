package activities

import (
	"elogika.vsb.cz/backend/modules/activities/handlers"
	"elogika.vsb.cz/backend/modules/auth/wrappers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("courses/:courseId/activities/:courseItemId", wrappers.WithUserDataRole(handlers.List))
	rg.GET("courses/:courseId/activities/:courseItemId/:termId", wrappers.WithUserDataRole(handlers.List))

	rg.GET("courses/:courseId/activities/available", wrappers.WithUserData(handlers.ListAvailable))
	rg.PUT("courses/:courseId/activities/prepare", wrappers.WithUserDataRole(handlers.ActivityInstancePrepare))
	rg.GET("courses/:courseId/activities/instance/:instanceId", wrappers.WithUserDataRole(handlers.ActivityInstanceGet))
	rg.PUT("courses/:courseId/activities/instance/:instanceId", wrappers.WithUserDataRole(handlers.ActivityInstanceSave))
	rg.DELETE("courses/:courseId/activities/instance/:instanceId", wrappers.WithUserDataRole(handlers.ActivityInstanceDelete))
}
