package handlers

import (
	"fmt"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/services/inbus"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Updated list of tutors
type ClassImportClassesResponse struct {
	Success bool `json:"success"`
}

// @Summary Removes tutor from class
// @Tags Classes
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} ClassImportStudentsResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/classes/{classId}/students/import [post]
func ImportClasses(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		any,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}
	if userRole != enums.CourseUserRoleAdmin && userRole != enums.CourseUserRoleGarant {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	courseService := services.NewCourseService(repositories.NewCourseRepository())
	course, err := courseService.GetCourseByID(initializers.DB, params.CourseID, userData.ID, userRole, nil, false, nil)
	if err != nil {
		return err
	}

	inbusClient := inbus.GetInbusClient()

	// Get set-up semester ID
	semester, err := inbusClient.GetSemesterFromDate(course.ImportOptions.Date)
	if err != nil {
		return err
	}

	// Get subject code
	subject, err := inbusClient.GetSubjectVersionFromcode(course.ImportOptions.Code)
	if err != nil {
		return err
	}

	// Get schedule activitity
	concreteActivities, err := inbusClient.GetConcreteActivities(subject.SubjectVersionId, &semester.SemesterId)
	if err != nil {
		return err
	}

	transaction := initializers.DB.Begin()

	// Already existing classes
	classRepo := repositories.ClassRepository{}
	existingClasses, _, err := classRepo.ListClasses(transaction, params.CourseID, userData.ID, nil, false, nil)
	if err != nil {
		return err
	}

	// Create or update classes
	for _, concreteActivity := range *concreteActivities {
		importCode := concreteActivity.ImportCode()
		ctype, form := concreteActivity.ClassType()
		day, _ := concreteActivity.Day()

		foundClass := GetClassIfExists(existingClasses, importCode, form, ctype)
		if foundClass == nil {
			foundClass = &models.Class{
				Version:      1,
				CourseID:     params.CourseID,
				StudentLimit: 40,
			}
		}

		foundClass.Name = concreteActivity.ClassName()
		foundClass.Room = concreteActivity.RoomName()
		foundClass.Type = ctype
		foundClass.StudyForm = form
		foundClass.TimeFrom = concreteActivity.BeginTime[:5]
		foundClass.TimeTo = concreteActivity.EndTime[:5]
		foundClass.Day = day
		foundClass.WeekParity = concreteActivity.WeekParity()
		foundClass.ImportOptions = models.ClassImportOptions{
			Code: concreteActivity.ImportCode(),
		}

		// TODO import students/tutors
		// foundClass.Students []ClassStudent ``
		// foundClass.Tutors   []ClassTutor   ``

		if foundClass.ID == 0 {
			existingClasses = append(existingClasses, foundClass)
		}

		if err := transaction.Save(&foundClass).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to create or update class",
			}
		}

		if concreteActivity.ConcreteActivityId == 312275 {
			utils.DebugPrintJSON(concreteActivity)
			utils.DebugPrintJSON(foundClass)
		}
	}

	// For each class import students
	existingClasses, _, err = classRepo.ListClasses(transaction, params.CourseID, userData.ID, nil, true, nil)
	if err != nil {
		return err
	}

	for _, class := range existingClasses {
		switch class.StudyForm {
		case enums.StudyFormFulltime:
			err = ImportByConcreteActivity(transaction, course, userRole, class)
			if err != nil {
				transaction.Rollback()
				return err
			}
		case enums.StudyFormCombined:
			err = ImportPartTimeStudents(transaction, course, userRole, class)
			if err != nil {
				transaction.Rollback()
				return err
			}
		default:
			panic(fmt.Sprintf("unexpected enums.ClassTypeEnum: %#v", class.Type))
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	c.JSON(200, ClassImportStudentsResponse{
		Success: true,
	})
	return nil
}

func GetClassIfExists(classes []*models.Class, code string, studyForm enums.StudyFormEnum, ctype enums.ClassTypeEnum) *models.Class {
	for _, class := range classes {
		if class.StudyForm == studyForm {
			switch studyForm {
			case enums.StudyFormCombined:
				if class.Type == ctype {
					return class
				}
			case enums.StudyFormFulltime:
				if class.ImportOptions.Code == code && class.Type == ctype {
					return class
				}
			default:
				panic(fmt.Sprintf("unexpected enums.StudyFormEnum: %#v", studyForm))
			}
		}
	}
	return nil
}
