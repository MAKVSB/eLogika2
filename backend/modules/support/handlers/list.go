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

type SupportTicketListResponse struct {
	Items      []dtos.SupportTicketListItemDTO `json:"items"`
	ItemsCount int64                           `json:"itemsCount"`
}

type SupportTicketListRequest struct {
	common.SearchRequest
}

// @Summary List all available support tickets in course
// @Tags SupportTickets
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body SupportTicketListRequest true "Ability to filter results"
// @Success 200 {object} SupportTicketListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/support [get]
func List(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, _, _, searchParams := utils.GetRequestDataWithSearch[
		any,
		any,
	](c, "search")
	if err != nil {
		return err
	}

	// TODO validate from here

	isAdmin := userRole == enums.CourseUserRoleAdmin && userData.Type == enums.UserTypeAdmin
	supportTicketService := services.NewSupportTicketService(repositories.NewSupportTicketRepository())
	tickets, ticketsCount, err := supportTicketService.ListSupportTickets(initializers.DB, userData.ID, isAdmin, nil, true, searchParams)
	if err != nil {
		return err
	}

	dtoList := make([]dtos.SupportTicketListItemDTO, len(tickets))
	for i, q := range tickets {
		dtoList[i] = dtos.SupportTicketListItemDTO{}.From(q)
	}

	c.JSON(200, SupportTicketListResponse{
		Items:      dtoList,
		ItemsCount: ticketsCount,
	})

	return nil
}
