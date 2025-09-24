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
	if userRole == enums.CourseUserRoleAdmin {
		return r.templateRepo.GetTemplateByIDAdmin(dbRef, courseID, templateID, userID, full, version)
	} else if userRole == enums.CourseUserRoleGarant {
		return r.templateRepo.GetTemplateByIDGarant(dbRef, courseID, templateID, userID, full, version)
	} else if userRole == enums.CourseUserRoleTutor {
		return r.templateRepo.GetTemplateByIDTutor(dbRef, courseID, templateID, userID, full, version)
	} else {
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
	if userRole == enums.CourseUserRoleAdmin {
		return r.templateRepo.ListTemplatesAdmin(dbRef, courseID, userID, full, searchParams)
	} else if userRole == enums.CourseUserRoleGarant {
		return r.templateRepo.ListTemplatesGarant(dbRef, courseID, userID, full, searchParams)
	} else if userRole == enums.CourseUserRoleTutor {
		return r.templateRepo.ListTemplatesTutor(dbRef, courseID, userID, full, searchParams)
	} else {
		return nil, 0, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}
