package services

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/repositories"
	"gorm.io/gorm"
)

type SupportTicketService struct {
	supportTicketRepo *repositories.SupportTicketRepository
}

func NewSupportTicketService(repo *repositories.SupportTicketRepository) *SupportTicketService {
	return &SupportTicketService{supportTicketRepo: repo}
}

func (r *SupportTicketService) GetSupportTicketByID(
	dbRef *gorm.DB,
	ticketID uint,
	userID uint,
	isAdmin bool,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
) (*models.SupportTicket, *common.ErrorResponse) {
	if isAdmin {
		return r.supportTicketRepo.GetSupportTicketByID(dbRef, ticketID, filters, full)
	} else {
		modifier := func(db *gorm.DB) *gorm.DB {
			if filters != nil {
				db = (*filters)(db)
			}
			return db.Where("created_by_id = ?", userID)
		}
		return r.supportTicketRepo.GetSupportTicketByID(dbRef, ticketID, &modifier, full)
	}
}

func (r *SupportTicketService) ListSupportTickets(
	dbRef *gorm.DB,
	userID uint,
	isAdmin bool,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.SupportTicket, int64, *common.ErrorResponse) {
	if isAdmin {
		return r.supportTicketRepo.ListSupportTickets(dbRef, filters, full, searchParams)
	} else {
		modifier := func(db *gorm.DB) *gorm.DB {
			if filters != nil {
				db = (*filters)(db)
			}
			return db.Where("created_by_id = ?", userID)
		}
		return r.supportTicketRepo.ListSupportTickets(dbRef, &modifier, full, searchParams)
	}
}
