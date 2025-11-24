package repositories

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"gorm.io/gorm"
)

type SupportTicketRepository struct{}

func NewSupportTicketRepository() *SupportTicketRepository {
	return &SupportTicketRepository{}
}

func (r *SupportTicketRepository) GetSupportTicketByID(
	dbRef *gorm.DB,
	supportTicketID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
) (*models.SupportTicket, *common.ErrorResponse) {

	query := dbRef.
		Where("id = ?", supportTicketID).
		Preload("CreatedBy").
		Preload("UpdatedBy")

	if full {
		query = query.
			Preload("Comments").
			Preload("Comments.CreatedBy")
	}

	if filters != nil {
		query = (*filters)(query)
	}

	var supportTicket *models.SupportTicket
	if err := query.
		Find(&supportTicket).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    401,
			Message: "Failed to load supportTicket",
		}
	}

	return supportTicket, nil
}

func (r *SupportTicketRepository) ListSupportTickets(
	dbRef *gorm.DB,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.SupportTicket, int64, *common.ErrorResponse) {
	var err *common.ErrorResponse
	query := dbRef.
		Model(models.SupportTicket{}).
		Preload("CreatedBy").
		Preload("UpdatedBy")

	if full {
		query = query.Preload("Comments")
	}

	if filters != nil {
		query = (*filters)(query)
	}

	// Apply filters, sorting, pagination
	if searchParams != nil {
		query, err = models.SupportTicket{}.ApplyFilters(query, searchParams.ColumnFilters, models.SupportTicket{}, map[string]interface{}{}, "")
		if err != nil {
			return nil, 0, err
		}
		query = models.SupportTicket{}.ApplySorting(query, searchParams.Sorting, "id DESC")
	}
	totalCount := models.SupportTicket{}.GetCount(query) // Gets count before pagination
	if searchParams != nil {
		query = models.SupportTicket{}.ApplyPagination(query, searchParams.Pagination)
	}

	var supportTickets []*models.SupportTicket
	if err := query.
		Find(&supportTickets).Error; err != nil {
		return nil, 0, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch supportTicket",
			Details: err.Error(),
		}
	}

	return supportTickets, totalCount, nil
}
