package handlers

import (
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"elogika.vsb.cz/backend/utils/tiptap"
	"github.com/gin-gonic/gin"
)

// @Description Request to insert new support ticket
type SupportTicketCommentInsertRequest struct {
	Content *models.TipTapContent `json:"content" binding:"required" ts_type:"JSONContent"`
}

// @Description Newly created support ticket
type SupportTicketCommentInsertResponse struct {
	Success bool `json:"success"`
}

// @Summary Create new ticket
// @Description Adds new ticket as a child of `parentId`
// @Tags SupportTicketComments
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body SupportTicketCommentInsertRequest true "New data for support ticket"
// @Success 200 {object} SupportTicketCommentInsertResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/support [post]
func SupportTicketCommentInsert(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			TicketID uint `uri:"ticketId" binding:"required"`
		},
		SupportTicketCommentInsertRequest,
	](c)
	if err != nil {
		return err
	}

	isAdmin := userRole == enums.CourseUserRoleAdmin && userData.Type == enums.UserTypeAdmin
	supportTicketService := services.NewSupportTicketService(repositories.NewSupportTicketRepository())
	ticket, err := supportTicketService.GetSupportTicketByID(initializers.DB, params.TicketID, userData.ID, isAdmin, nil, false)
	if err != nil {
		return err
	}

	transaction := initializers.DB.Begin()

	comment := &models.SupportTicketComment{
		Content:         reqData.Content,
		CreatedByID:     userData.ID,
		SupportTicketID: ticket.ID,
	}

	if err := transaction.Save(&comment).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to save ticket comment",
			Details: err.Error(),
		}
	}

	err = tiptap.FindAndSaveRelations(transaction, userData.ID, reqData.Content, &comment, "ContentFiles")
	if err != nil {
		transaction.Rollback()
		return err
	}

	if err := transaction.Save(&comment).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update support ticket",
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	c.JSON(200, SupportTicketCommentInsertResponse{
		Success: true,
	})

	return nil
}
