package auth

import (
	"elogika.vsb.cz/backend/modules/auth/handlers"
	"elogika.vsb.cz/backend/modules/auth/wrappers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutesUnauth(rg *gin.RouterGroup) {
	rg.POST("auth/login", handlers.Login)
	rg.POST("auth/login/sso", handlers.SSOLogin)
	rg.POST("auth/login/sso/callback", handlers.SSOLoginCallback)
	rg.POST("auth/refresh", handlers.Refresh)
}

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("auth/logout", wrappers.WithToken(handlers.Logout))
	rg.POST("auth/logout/all", wrappers.WithToken(handlers.LogoutAll))
}
