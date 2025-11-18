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
type SupportTicketUpdateRequest struct {
	Name    string                `json:"name"`
	Content *models.TipTapContent `json:"content" binding:"required" ts_type:"JSONContent"`
	Solved  bool                  `json:"solved"`
	URL     string                `json:"url"`
}

// @Description Updated support ticket
type SupportTicketUpdateResponse struct {
	Data dtos.SupportTicketDTO `json:"data"`
}

type CourseSupportTicketUri struct {
	TicketID uint `uri:"ticketId" binding:"required"`
}

// @Summary Update support ticket
// @Description Updates support ticket content
// @Tags SupportTickets
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param ticketId path int true "ID of the edited support ticket"
// @Param body body SupportTicketUpdateRequest true "New data for support ticket"
// @Success 200 {object} SupportTicketUpdateResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/support/{ticketId} [put]
func SupportTicketUpdate(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			TicketID uint `uri:"ticketId" binding:"required"`
		},
		SupportTicketUpdateRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	transaction := initializers.DB.Begin()

	isAdmin := userRole == enums.CourseUserRoleAdmin && userData.Type == enums.UserTypeAdmin
	supportTicketService := services.NewSupportTicketService(repositories.NewSupportTicketRepository())

	ticket, err := supportTicketService.GetSupportTicketByID(transaction, params.TicketID, userData.ID, isAdmin, nil, false)
	if err != nil {
		transaction.Rollback()
		return err
	}
	// Partially modify data
	ticket.Name = reqData.Name
	err = tiptap.FindAndSaveRelations(transaction, userData.ID, reqData.Content, &ticket, "ContentFiles")
	if err != nil {
		return err
	}
	ticket.Content = reqData.Content
	ticket.Solved = reqData.Solved
	ticket.URL = reqData.URL

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

	ticket, err = supportTicketService.GetSupportTicketByID(transaction, params.TicketID, userData.ID, isAdmin, nil, false)
	if err != nil {
		return err
	}

	c.JSON(200, SupportTicketUpdateResponse{
		Data: dtos.SupportTicketDTO{}.From(ticket, isAdmin),
	})
	return nil
}
