package users

import (
	"elogika.vsb.cz/backend/modules/auth/wrappers"
	"elogika.vsb.cz/backend/modules/users/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("users", wrappers.WithUserData(handlers.List))

	rg.GET("users/self", wrappers.WithUserDataRole(handlers.GetSelf))
	rg.PUT("users/self", wrappers.WithUserDataRole(handlers.UserUpdateSelf))

	rg.GET("users/:userId", wrappers.WithUserData(handlers.GetByID))
	rg.POST("users", wrappers.WithUserData(handlers.UserInsert))
	rg.PUT("users/:userId", wrappers.WithUserData(handlers.UserUpdate))
}
