package services

import (
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"gorm.io/gorm"
)

type CourseService struct {
	courseRepo *repositories.CourseRepository
}

func NewCourseService(repo *repositories.CourseRepository) *CourseService {
	return &CourseService{courseRepo: repo}
}

func (r *CourseService) GetCourseByID(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	userRole enums.CourseUserRoleEnum,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
) (*models.Course, *common.ErrorResponse) {
	if userRole == enums.CourseUserRoleAdmin || userRole == enums.CourseUserRoleGarant || userRole == enums.CourseUserRoleTutor {
		return r.courseRepo.GetCourseByID(dbRef, courseID, userID, nil, full, version)
	} else if userRole == enums.CourseUserRoleStudent {

		var courseUser *models.CourseUser
		if err := initializers.DB.
			Where("course_id = ?", courseID).
			Where("user_id = ?", userID).
			First(&courseUser).Error; err != nil {
			return nil, &common.ErrorResponse{
				Code:    500,
				Message: "Failed to commit changes",
			}
		}
		return r.courseRepo.GetCourseByID(dbRef, courseID, userID, nil, full, version)
	} else {
		modifier := func(db *gorm.DB) *gorm.DB {
			return db.Where("public = ?", true)
		}
		return r.courseRepo.GetCourseByID(dbRef, courseID, userID, &modifier, full, version)
	}
}

func (s *CourseService) GetUserCourses(
	userID uint,
	userType *enums.UserTypeEnum,
) ([]*models.CourseUser, *common.ErrorResponse) {
	return s.courseRepo.GetUserCoursesByUserID(initializers.DB, userID, userType)
}
