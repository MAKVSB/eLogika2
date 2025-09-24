package services

import (
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"gorm.io/gorm"
)

type TermService struct {
	termRepo *repositories.TermRepository
}

func NewTermService(repo *repositories.TermRepository) *TermService {
	return &TermService{termRepo: repo}
}

func (r *TermService) GetTermByID(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	termID uint,
	userID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
) (*models.Term, *common.ErrorResponse) {
	if userRole == enums.CourseUserRoleAdmin {
		// Can user see the course item
		cis := services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
		courseItem, err := cis.GetCourseItemByID(dbRef, courseID, courseItemID, userID, userRole, nil, false, nil)
		if err != nil {
			return nil, err
		}
		if !courseItem.Editable {
			return nil, &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permissions",
			}
		}

		return r.termRepo.GetTermByIDAdmin(initializers.DB, courseID, courseItemID, termID, userID, full, version)
	} else if userRole == enums.CourseUserRoleGarant {
		// Can user see the course item
		cis := services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
		courseItem, err := cis.GetCourseItemByID(dbRef, courseID, courseItemID, userID, userRole, nil, false, nil)
		if err != nil {
			return nil, err
		}
		if !courseItem.Editable {
			return nil, &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permissions",
			}
		}

		return r.termRepo.GetTermByIDGarant(initializers.DB, courseID, courseItemID, termID, userID, full, version)
	} else if userRole == enums.CourseUserRoleTutor {
		// Can user see the course item
		cis := services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
		courseItem, err := cis.GetCourseItemByID(dbRef, courseID, courseItemID, userID, userRole, nil, false, nil)
		if err != nil {
			return nil, err
		}
		if !courseItem.Editable {
			return nil, &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permissions",
			}
		}

		return r.termRepo.GetTermByIDTutor(initializers.DB, courseID, courseItemID, termID, userID, full, version)
	} else {
		return nil, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}

func (r *TermService) ListTerms(
	dbRef *gorm.DB,
	courseID uint,
	courseItemID uint,
	userID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.Term, int64, *common.ErrorResponse) {
	if userRole == enums.CourseUserRoleAdmin {
		// Can user see the course item
		cis := services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
		courseItem, err := cis.GetCourseItemByID(dbRef, courseID, courseItemID, userID, userRole, nil, false, nil)
		if err != nil {
			return nil, 0, err
		}
		if !courseItem.Editable {
			return nil, 0, &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permissions",
			}
		}

		return r.termRepo.ListTermsAdmin(dbRef, courseID, courseItemID, userID, full, searchParams)
	} else if userRole == enums.CourseUserRoleGarant {
		// Can user see the course item
		cis := services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
		courseItem, err := cis.GetCourseItemByID(dbRef, courseID, courseItemID, userID, userRole, nil, false, nil)
		if err != nil {
			return nil, 0, err
		}
		if !courseItem.Editable {
			return nil, 0, &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permissions",
			}
		}

		return r.termRepo.ListTermsGarant(dbRef, courseID, courseItemID, userID, full, searchParams)
	} else if userRole == enums.CourseUserRoleTutor {
		// Can user see the course item
		cis := services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
		courseItem, err := cis.GetCourseItemByID(dbRef, courseID, courseItemID, userID, userRole, nil, false, nil)
		if err != nil {
			return nil, 0, err
		}
		if !courseItem.Editable {
			return nil, 0, &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permissions",
			}
		}

		return r.termRepo.ListTermsTutor(dbRef, courseID, courseItemID, userID, full, searchParams)
	} else {
		// TODO student will probably also have access to this endpoint
		return nil, 0, &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
}

func (r *TermService) ListJoinedStudents(
	dbRef *gorm.DB,
	termID uint,
	userID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	searchParams *common.SearchRequest,
) ([]*models.UserTerm, int64, *common.ErrorResponse) {
	return r.termRepo.ListJoinedStudents(dbRef, termID, nil, false, nil)
}
