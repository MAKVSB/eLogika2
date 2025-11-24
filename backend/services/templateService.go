package services

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"gorm.io/gorm"
)

type TemplateService struct {
	templateRepo *repositories.TemplateRepository
}

func NewTemplateService(repo *repositories.TemplateRepository) *TemplateService {
	return &TemplateService{templateRepo: repo}
}

func (r *TemplateService) GetTemplateByID(
	dbRef *gorm.DB,
	courseID uint,
	templateID uint,
	userID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
) (*models.Template, *common.ErrorResponse) {
	switch userRole {
	case enums.CourseUserRoleAdmin:
		return r.templateRepo.GetTemplateByID(dbRef, courseID, templateID, userID, filters, full, version)
	case enums.CourseUserRoleGarant:
		modifier := func(db *gorm.DB) *gorm.DB {
			if filters != nil {
				db = (*filters)(db)
			}
			return db.Where("managed_by = ?", enums.CourseUserRoleGarant)
		}
		return r.templateRepo.GetTemplateByID(dbRef, courseID, templateID, userID, &modifier, full, version)
	case enums.CourseUserRoleTutor:
		modifier := func(db *gorm.DB) *gorm.DB {
			if filters != nil {
				db = (*filters)(db)
			}
			return db.Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userID)
		}
		return r.templateRepo.GetTemplateByID(dbRef, courseID, templateID, userID, &modifier, full, version)
	default:
		return nil, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}

func (r *TemplateService) ListTemplates(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Template, int64, *common.ErrorResponse) {
	switch userRole {
	case enums.CourseUserRoleAdmin:
		return r.templateRepo.ListTemplates(dbRef, courseID, userID, filters, full, searchParams)
	case enums.CourseUserRoleGarant:
		modifier := func(db *gorm.DB) *gorm.DB {
			if filters != nil {
				db = (*filters)(db)
			}
			return db.Where("managed_by = ?", enums.CourseUserRoleGarant)
		}
		return r.templateRepo.ListTemplates(dbRef, courseID, userID, &modifier, full, searchParams)
	case enums.CourseUserRoleTutor:
		modifier := func(db *gorm.DB) *gorm.DB {
			if filters != nil {
				db = (*filters)(db)
			}
			return db.Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userID)
		}
		return r.templateRepo.ListTemplates(dbRef, courseID, userID, &modifier, full, searchParams)
	default:
		return nil, 0, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}
