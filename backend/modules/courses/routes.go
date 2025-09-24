package courses

import (
	"elogika.vsb.cz/backend/modules/auth/wrappers"
	"elogika.vsb.cz/backend/modules/courses/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("courses", wrappers.WithUserDataRole(handlers.CourseInsert))
	rg.PUT("courses/:courseId", wrappers.WithUserDataRole(handlers.Update))
	rg.GET("courses", wrappers.WithUserDataRole(handlers.List))
	rg.GET("courses/user", wrappers.WithUserDataRole(handlers.ListUserCourses))
	rg.GET("courses/:courseId", wrappers.WithUserDataRole(handlers.CourseGetByID))

	rg.GET("courses/:courseId/users", wrappers.WithUserDataRole(handlers.ListCourseUsers))
	rg.POST("courses/:courseId/users", wrappers.WithUserDataRole(handlers.AddCourseUser))
	rg.DELETE("courses/:courseId/users", wrappers.WithUserDataRole(handlers.RemoveCourseUser))
}
