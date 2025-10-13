package tests

import (
	"elogika.vsb.cz/backend/modules/auth/wrappers"
	"elogika.vsb.cz/backend/modules/tests/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("courses/:courseId/tests/:courseItemId", wrappers.WithUserDataRole(handlers.List))
	rg.GET("courses/:courseId/tests/:courseItemId/:termId", wrappers.WithUserDataRole(handlers.List))
	rg.POST("courses/:courseId/tests/:courseItemId/:termId/generate", wrappers.WithUserDataRole(handlers.Generate))
	rg.GET("courses/:courseId/tests/:courseItemId/instances/:testId", wrappers.WithUserDataRole(handlers.ListInstance))
	rg.DELETE("courses/:courseId/tests/:courseItemId/instances/:testId", wrappers.WithUserDataRole(handlers.TestDelete))
	rg.POST("courses/:courseId/tests/:courseItemId/instances/:testId/create", wrappers.WithUserDataRole(handlers.CreateInstance))

	rg.PUT("courses/:courseId/tests/:courseItemId/instance/:instanceId/tutorsave", wrappers.WithUserDataRole(handlers.TestInstanceTutorSave))
	rg.GET("courses/:courseId/tests/:courseItemId/instance/:instanceId", wrappers.WithUserDataRole(handlers.TestInstanceTutorGet))
	rg.GET("courses/:courseId/tests/:courseItemId/instance/:instanceId/telemetry", wrappers.WithUserDataRole(handlers.TestInstanceGetTelemetry))
	rg.DELETE("courses/:courseId/tests/:courseItemId/instance/:instanceId", wrappers.WithUserDataRole(handlers.TestInstanceDelete))

	rg.GET("courses/:courseId/tests/available", wrappers.WithUserData(handlers.ListAvailable))
	rg.POST("courses/:courseId/tests/prepare", wrappers.WithUserData(handlers.TestInstancePrepare))
	rg.POST("courses/:courseId/tests/evaluate", wrappers.WithUserDataRole(handlers.TestEvaluate))

	rg.GET("tests/:instanceId", wrappers.WithUserDataRole(handlers.TestInstanceGet))
	rg.PUT("tests/:instanceId/start", wrappers.WithUserDataRole(handlers.TestInstanceStart))
	rg.PUT("tests/:instanceId/save", wrappers.WithUserDataRole(handlers.TestInstanceSave))
	rg.PUT("tests/:instanceId/finish", wrappers.WithUserDataRole(handlers.TestInstanceFinish))
	rg.POST("tests/:instanceId/telemetry", wrappers.WithUserDataRole(handlers.TestInstanceTelemetry))
}
