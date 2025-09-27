package classes

import (
	"elogika.vsb.cz/backend/modules/auth/wrappers"
	"elogika.vsb.cz/backend/modules/classes/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("courses/:courseId/classes", wrappers.WithUserDataRole(handlers.List))
	// rg.POST("courses/:courseId/classes/import", wrappers.WithUserDataRole(handlers.ImportClasses))
	rg.POST("courses/:courseId/classes", wrappers.WithUserDataRole(handlers.ClassInsert))
	rg.PUT("courses/:courseId/classes/:classId", wrappers.WithUserDataRole(handlers.ClassUpdate))
	rg.GET("courses/:courseId/classes/:classId", wrappers.WithUserDataRole(handlers.ClassGetByID))
	rg.DELETE("courses/:courseId/classes/:classId", wrappers.WithUserDataRole(handlers.ClassDelete))

	rg.GET("courses/:courseId/classes/:classId/students", wrappers.WithUserDataRole(handlers.ListStudents))
	rg.POST("courses/:courseId/classes/:classId/students/import", wrappers.WithUserDataRole(handlers.ImportStudents))
	rg.POST("courses/:courseId/classes/:classId/students", wrappers.WithUserDataRole(handlers.AddStudent))
	rg.DELETE("courses/:courseId/classes/:classId/students", wrappers.WithUserDataRole(handlers.RemoveStudent))

	rg.GET("courses/:courseId/classes/:classId/tutors", wrappers.WithUserDataRole(handlers.ListTutors))
	rg.POST("courses/:courseId/classes/:classId/tutors", wrappers.WithUserDataRole(handlers.AddTutor))
	rg.DELETE("courses/:courseId/classes/:classId/tutors", wrappers.WithUserDataRole(handlers.RemoveTutor))
}
