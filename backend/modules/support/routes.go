package support

import (
	"elogika.vsb.cz/backend/modules/auth/wrappers"
	"elogika.vsb.cz/backend/modules/support/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("support", wrappers.WithUserDataRole(handlers.SupportTicketInsert))
	rg.PUT("support/:ticketId", wrappers.WithUserDataRole(handlers.SupportTicketUpdate))
	rg.PUT("support/:ticketId/comment", wrappers.WithUserDataRole(handlers.SupportTicketCommentInsert))
	rg.GET("support/:ticketId", wrappers.WithUserDataRole(handlers.SupportTicketGetByID))
	rg.GET("support", wrappers.WithUserDataRole(handlers.List))
}
