package services

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"gorm.io/gorm"
)

type ChapterService struct {
	chapterRepo *repositories.ChapterRepository
}

func NewChapterService(repo *repositories.ChapterRepository) *ChapterService {
	return &ChapterService{chapterRepo: repo}
}

func (r *ChapterService) GetChapterByID(
	dbRef *gorm.DB,
	courseID uint,
	chapterID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
) (*models.Chapter, *common.ErrorResponse) {
	switch userRole {
	case enums.CourseUserRoleAdmin, enums.CourseUserRoleGarant, enums.CourseUserRoleTutor:
		return r.chapterRepo.GetChapterByID(dbRef, courseID, chapterID, filters, full, version, false)
	case enums.CourseUserRoleStudent:
		return r.chapterRepo.GetChapterByID(dbRef, courseID, chapterID, filters, full, version, true)
	default:
		return nil, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}

func (r *ChapterService) ListChapters(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Chapter, int64, *common.ErrorResponse) {
	switch userRole {
	case enums.CourseUserRoleAdmin, enums.CourseUserRoleGarant, enums.CourseUserRoleTutor:
		return r.chapterRepo.ListChapters(dbRef, courseID, filters, full, searchParams, false)
	case enums.CourseUserRoleStudent:
		return r.chapterRepo.ListChapters(dbRef, courseID, filters, full, searchParams, true)
	default:
		return nil, 0, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}
