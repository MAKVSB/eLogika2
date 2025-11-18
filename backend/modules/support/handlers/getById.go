package handlers

import (
	"elogika.vsb.cz/backend/initializers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/support/dtos"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Newly created support ticket
type SupportTicketGetByIdResponse struct {
	Data dtos.SupportTicketDTO `json:"data"`
}

// @Summary Get support ticket by id
// @Tags SupportTickets
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param support ticketId path int true "ID of the edited support ticket"
// @Success 200 {object} SupportTicketGetByIdResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/support/{ticketId} [get]
func SupportTicketGetByID(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			TicketID uint `uri:"ticketId" binding:"required"`
		},
		any,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	isAdmin := userRole == enums.CourseUserRoleAdmin && userData.Type == enums.UserTypeAdmin
	supportTicketService := services.NewSupportTicketService(repositories.NewSupportTicketRepository())
	ticket, err := supportTicketService.GetSupportTicketByID(initializers.DB, params.TicketID, userData.ID, isAdmin, nil, true)
	if err != nil {
		return err
	}

	c.JSON(200, SupportTicketInsertResponse{
		Data: dtos.SupportTicketDTO{}.From(ticket, isAdmin),
	})

	return nil
}
