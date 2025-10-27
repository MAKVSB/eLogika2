package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	authenums "elogika.vsb.cz/backend/modules/auth/enums"
	"elogika.vsb.cz/backend/modules/auth/tokens"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/utils"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// @Description Request to insert user into class as tutor
type RemoveCourseUserRequest struct {
	UserID uint                      `json:"userId" binding:"required"`
	Role   *enums.CourseUserRoleEnum `json:"role"`
}

// @Description Updated list of students
type RemoveCourseUserResponse struct {
	Success bool `json:"success"`
}

// @Summary Add user to course
// @Tags Courses
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body RemoveCourseUserRequest true "New data for class"
// @Success 200 {object} RemoveCourseUserResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/users [delete]
func RemoveCourseUser(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		RemoveCourseUserRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}
	// If not admin, garant, or tutor
	if userRole != enums.CourseUserRoleAdmin {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	var courseUser models.CourseUser
	if err := initializers.DB.
		Where("user_id = ?", reqData.UserID).
		Where("course_id = ?", params.CourseID).
		First(&courseUser).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch course user",
		}
	}

	if reqData.Role != nil {
		var newRoles []enums.CourseUserRoleEnum

		switch *reqData.Role {
		case enums.CourseUserRoleStudent:
			err := checkStudentNoClasses(initializers.DB, reqData.UserID, params.CourseID)
			if err != nil {
				return err
			}
		case enums.CourseUserRoleTutor:
			err := checkTutorNoClasses(initializers.DB, reqData.UserID, params.CourseID)
			if err != nil {
				return err
			}
		}

		for _, role := range courseUser.Roles {
			if role != *reqData.Role {
				newRoles = append(newRoles, role)
			}
		}

		if len(newRoles) == 0 {
			if err := initializers.DB.
				Where("user_id = ?", reqData.UserID).
				Where("course_id = ?", params.CourseID).
				Delete(&models.CourseUser{}).Error; err != nil {
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to remove course user",
					Details: err.Error(),
				}
			}
		} else {
			courseUser.Roles = newRoles
			if err := initializers.DB.
				Save(&courseUser).Error; err != nil {
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to update course user",
					Details: err.Error(),
				}
			}
		}
	} else {
		// Called API to remove whole user

		// Check if user is not in class as student
		err := checkStudentNoClasses(initializers.DB, reqData.UserID, params.CourseID)
		if err != nil {
			return err
		}

		// Check if user is not in class as tutor
		err = checkTutorNoClasses(initializers.DB, reqData.UserID, params.CourseID)
		if err != nil {
			return err
		}

		if err := initializers.DB.
			Where("user_id = ?", reqData.UserID).
			Where("course_id = ?", params.CourseID).
			Delete(&models.CourseUser{}).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to remove course user",
				Details: err.Error(),
			}
		}
	}

	accessToken := tokens.AccessToken{}
	accessToken.UserID = reqData.UserID
	accessToken.TokenType = authenums.JWTTokenTypeAccess
	accessToken.InvalidateByUser()

	c.JSON(200, RemoveCourseUserResponse{
		Success: true,
	})

	return nil
}

func checkTutorNoClasses(dbRef *gorm.DB, userId uint, courseId uint) *common.ErrorResponse {
	var tutorClasses []*models.ClassTutor
	if err := dbRef.
		Where("user_id = ?", userId).
		InnerJoins("Class", initializers.DB.Where("course_id = ?", courseId)).
		Find(&tutorClasses).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to get user classes",
			Details: err.Error(),
		}
	}

	if len(tutorClasses) != 0 {
		errResources := make([]common.ErrorResources, len(tutorClasses))
		for i, studentClasses := range tutorClasses {
			errResources[i] = common.ErrorResources{
				ResourceType: "class",
				ResourceID:   studentClasses.ClassID,
			}
		}

		return &common.ErrorResponse{
			Code:      409,
			Message:   "Failed to remove user from course",
			Details:   "User is in class",
			Resources: errResources,
		}
	}
	return nil
}

func checkStudentNoClasses(dbRef *gorm.DB, userId uint, courseId uint) *common.ErrorResponse {
	var tutorClasses []*models.ClassStudent
	if err := dbRef.
		Where("user_id = ?", userId).
		InnerJoins("Class", initializers.DB.Where("course_id = ?", courseId)).
		Find(&tutorClasses).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to get user classes",
			Details: err.Error(),
		}
	}

	if len(tutorClasses) != 0 {
		errResources := make([]common.ErrorResources, len(tutorClasses))
		for i, studentClasses := range tutorClasses {
			errResources[i] = common.ErrorResources{
				ResourceType: "class",
				ResourceID:   studentClasses.ClassID,
			}
		}

		return &common.ErrorResponse{
			Code:      409,
			Message:   "Failed to remove user from course",
			Details:   "User is in class",
			Resources: errResources,
		}
	}
	return nil
}
