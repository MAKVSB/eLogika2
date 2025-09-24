package repositories

import (
	"encoding/json"
	"strconv"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/utils"
	"gorm.io/gorm"
)

type CourseRepository struct{}

func NewCourseRepository() *CourseRepository {
	return &CourseRepository{}
}

func (r *CourseRepository) GetCourseByID(
	dbRef *gorm.DB,
	courseID uint,
	userID uint,
	filters *(func(*gorm.DB) *gorm.DB),
	full bool,
	version *uint,
) (*models.Course, *common.ErrorResponse) {
	query := dbRef

	if filters != nil {
		query = (*filters)(query)
	}

	if full {
	}

	var course *models.Course
	if err := query.
		First(&course, courseID).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch question",
			Details: err.Error(),
		}
	}

	if version != nil {
		if course.Version != *version {
			return nil, &common.ErrorResponse{
				Code:    409,
				Message: "Version mismatched",
				Details: strconv.Itoa(int(*version)) + " " + strconv.Itoa(int(course.Version)),
			}
		}
	}

	return course, nil
}

func (r *CourseRepository) GetUserCoursesByUserID(
	dbRef *gorm.DB,
	userID uint,
	userType *enums.UserTypeEnum,
) ([]*models.CourseUser, *common.ErrorResponse) {
	if userType != nil && *userType == enums.UserTypeAdmin {
		var courses []models.Course
		if err := dbRef.
			Find(&courses).Error; err != nil {
			return nil, &common.ErrorResponse{
				Code:    401,
				Message: "Failed to load courses",
			}
		}

		userCourses := make([]*models.CourseUser, len(courses))
		for i, c := range courses {
			userCourses[i] = &models.CourseUser{
				CourseID: c.ID,
				UserID:   userID,
				Roles:    enums.CourseUserRoleEnumAll,
				Course:   &c,
			}
		}

		return userCourses, nil
	} else {
		var userCourses []*models.CourseUser
		if err := dbRef.
			Preload("Course").
			Where("user_id = ?", userID).Find(&userCourses).Error; err != nil {
			return nil, &common.ErrorResponse{
				Code:    401,
				Message: "Failed to load courses",
			}
		}
		return userCourses, nil
	}
}

func (r *CourseRepository) GetCourseGarantsIds(
	dbRef *gorm.DB,
	courseId uint,
) ([]uint, *common.ErrorResponse) {

	var garantIDs []uint

	if err := dbRef.
		Model(&models.CourseUser{}).
		Preload("Course").
		Where("roles like '%GARANT%'").
		Pluck("user_id", &garantIDs).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    401,
			Message: "Failed to load garants",
			Details: err.Error(),
		}
	}
	return garantIDs, nil
}

func (r *CourseRepository) SyncFiles(
	dbRef *gorm.DB,
	content json.RawMessage,
	course *models.Course,
) *common.ErrorResponse {
	var files []*models.File
	if err := dbRef.Where("id IN ?", utils.GetFilesInsideContent(content)).Find(&files).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to load files",
		}
	}

	course.ContentFiles = files

	if err := dbRef.Model(&course).Association("ContentFiles").Replace(&course.ContentFiles); err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update files",
		}
	}
	return nil
}
