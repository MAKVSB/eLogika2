package handlers

import (
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/support/dtos"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"elogika.vsb.cz/backend/utils/tiptap"
	"github.com/gin-gonic/gin"
)

// @Description Request to insert new support ticket
type SupportTicketInsertRequest struct {
	Name    string                `json:"name"`
	Content *models.TipTapContent `json:"content" binding:"required" ts_type:"JSONContent"`
	Solved  bool                  `json:"solved"`
	URL     string                `json:"url"`
}

// @Description Newly created support ticket
type SupportTicketInsertResponse struct {
	Data dtos.SupportTicketDTO `json:"data"`
}

// @Summary Create new ticket
// @Description Adds new ticket as a child of `parentId`
// @Tags SupportTickets
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body SupportTicketInsertRequest true "New data for support ticket"
// @Success 200 {object} SupportTicketInsertResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/support [post]
func SupportTicketInsert(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, _, reqData := utils.GetRequestData[
		any,
		SupportTicketInsertRequest,
	](c)
	if err != nil {
		return err
	}

	transaction := initializers.DB.Begin()

	ticket := &models.SupportTicket{
		Name:        reqData.Name,
		Content:     reqData.Content,
		Solved:      reqData.Solved,
		CreatedByID: userData.ID,
		URL:         reqData.URL,
	}

	if err := transaction.Save(&ticket).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to save ticket",
		}
	}

	err = tiptap.FindAndSaveRelations(transaction, userData.ID, reqData.Content, &ticket, "ContentFiles")
	if err != nil {
		transaction.Rollback()
		return err
	}

	if err := transaction.Save(&ticket).Error; err != nil {
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

	isAdmin := userRole == enums.CourseUserRoleAdmin && userData.Type == enums.UserTypeAdmin
	supportTicketService := services.NewSupportTicketService(repositories.NewSupportTicketRepository())
	ticket, err = supportTicketService.GetSupportTicketByID(initializers.DB, ticket.ID, userData.ID, isAdmin, nil, true)
	if err != nil {
		return err
	}

	c.JSON(200, SupportTicketInsertResponse{
		Data: dtos.SupportTicketDTO{}.From(ticket, isAdmin),
	})

	return nil
}
