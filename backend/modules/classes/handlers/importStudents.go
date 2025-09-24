package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/classes/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services/inbus"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Updated list of tutors
type ClassImportStudentsResponse struct {
	Students []dtos.ClassUserDTO `json:"students"`
}

// @Summary Removes tutor from class
// @Tags Classes
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} RemoveStudentResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/ [post] // TODO
func ImportStudents(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
			ClassID  uint `uri:"classId"  binding:"required"`
		},
		any,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		return err
	}

	classRepo := repositories.NewClassRepository()
	var class *models.Class
	if userRole == enums.CourseUserRoleAdmin {
		class, err = classRepo.GetClassByIDAdmin(initializers.DB, params.CourseID, params.ClassID, userData.ID, false, nil)
	} else if userRole == enums.CourseUserRoleGarant {
		// Permission to update every class
		class, err = classRepo.GetClassByIDGarant(initializers.DB, params.CourseID, params.ClassID, userData.ID, false, nil)
	} else if userRole == enums.CourseUserRoleTutor {
		// Can only update his own class
		class, err = classRepo.GetClassByIDTutor(initializers.DB, params.CourseID, params.ClassID, userData.ID, false, nil)
	} else {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
	if err != nil {
		return err
	}

	var course *models.Course
	if err := initializers.DB.
		Find(&course, params.CourseID).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to remove student from class",
		}
	}

	var courseClasses []*models.Class
	if err := initializers.DB.
		Where("course_id = ?", params.CourseID).
		Find(&courseClasses).
		Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to remove student from class",
		}
	}

	courseClassesIDs := make([]uint, len(courseClasses))
	for _, courseClass := range courseClasses {
		courseClassesIDs = append(courseClassesIDs, courseClass.ID)
	}

	utils.DebugPrintJSON(course.ImportOptions)

	inbusClient := inbus.GetInbusClient()

	// Get set-up semester ID
	semester, err := inbusClient.GetSemesterFromDate(course.ImportOptions.Date)
	if err != nil {
		return err
	}
	if semester == nil {
		return &common.ErrorResponse{
			Code:    0,
			Message: "Failed load information from inbus",
			Details: "Semester not found",
		}
	}

	// Get subject code
	subject, err := inbusClient.GetSubjectVersionFromcode(course.ImportOptions.Code)
	if err != nil {
		return err
	}
	if subject == nil {
		return &common.ErrorResponse{
			Code:    0,
			Message: "Failed load information from inbus",
			Details: "Subject not found",
		}
	}

	concreteActivities, err := inbusClient.GetConcreteActivities((*subject)[0].SubjectVersionId, &semester.SemesterId)
	if err != nil {
		return err
	}
	if concreteActivities == nil || len(*concreteActivities) == 0 {
		return &common.ErrorResponse{
			Code:    0,
			Message: "Failed load information from inbus",
			Details: "Schedule is missing",
		}
	}

	var concreteActivityIds []uint
	for _, concreteActivity := range *concreteActivities {
		// Day of the week
		weekDayId := uint(0)
		switch class.Day {
		case enums.WeekDayMonday:
			weekDayId = 1
		case enums.WeekDayTuesday:
			weekDayId = 2
		case enums.WeekDayWednesday:
			weekDayId = 3
		case enums.WeekDayThursday:
			weekDayId = 4
		case enums.WeekDayFriday:
			weekDayId = 5
		case enums.WeekDaySaturday:
			weekDayId = 6
		case enums.WeekDaySunday:
			weekDayId = 7
		default:
			panic(fmt.Sprintf("unexpected enums.WeekDayEnum: %#v", class.Day))
		}
		if concreteActivity.WeekDayId != weekDayId {
			continue
		}

		// Room number
		if !strings.Contains(concreteActivity.RoomFullcodes, class.Room) {
			continue
		}

		// Time from
		if class.TimeFrom+":00" != concreteActivity.BeginTime {
			continue
		}
		// Time to
		if class.TimeTo+":00" != concreteActivity.EndTime {
			continue
		}

		// if class.StudyForm == enums.StudyFormFulltime {
		// 	// Class type
		// 	classTypeId := "P"
		// 	switch class.Type {
		// 	case enums.ClassTypeC:
		// 		classTypeId = "C"
		// 	case enums.ClassTypeP:
		// 		classTypeId = "P"
		// 	default:
		// 		panic(fmt.Sprintf("unexpected enums.ClassTypeEnum: %#v", class.Type))
		// 	}
		// 	if classTypeId != concreteActivity.EducationTypeAbbrev {
		// 		continue
		// 	}
		// }
		concreteActivityIds = append(concreteActivityIds, concreteActivity.ConcreteActivityId)
	}

	if len(concreteActivityIds) == 0 {
		return &common.ErrorResponse{
			Code:    409,
			Message: "Any of schedule items not matched",
		}
	}

	classAssignSkipped := make([]models.CourseUser, 0)

	transaction := initializers.DB.Begin()

	for _, concreteActivityId := range concreteActivityIds {

		concreteActivityStudents, err := inbusClient.GetConcreteActivityStudents(concreteActivityId)
		if err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    409,
				Message: "Failed to get students assignned to schedule item",
			}
		}

		for _, activityStudent := range *concreteActivityStudents {
			// If user does in exist in entire system, create. If does, update
			var user models.User
			{
				if err := transaction.
					Where("identity_provider_id = ?", activityStudent.PersonId).
					Find(&user).Error; err != nil {
					transaction.Rollback()
					return &common.ErrorResponse{
						Code:    500,
						Message: "Failed to remove student from class",
					}
				}

				user.Version = user.Version + 1
				user.Username = activityStudent.Login
				user.FirstName = activityStudent.FirstName
				user.FamilyName = activityStudent.SecondName
				user.Email = activityStudent.Login + "@vsb.cz"
				user.IdentityProvider = enums.IdentityProviderVSB
				user.IdentityProviderID = strconv.Itoa(int(activityStudent.PersonId))
				if user.ID == 0 {
					user.Notification = models.UserNotification{
						Discord: models.NotificationDiscord{
							Level: models.NotificationLevel{
								Results:  true,
								Messages: true,
								Terms:    true,
							},
							UserID: "",
						},
						Email: models.NotificationEmail{
							Level: models.NotificationLevel{
								Results:  true,
								Messages: true,
								Terms:    true,
							},
						},
						Push: models.NotificationPush{
							Level: models.NotificationLevel{
								Results:  true,
								Messages: true,
								Terms:    true,
							},
							Token: "",
						},
					}
					user.Type = enums.UserTypeNormal
				}

				if err := transaction.
					Save(&user).Error; err != nil {
					transaction.Rollback()
					return &common.ErrorResponse{
						Code:    500,
						Message: "Failed to save new user",
					}
				}
			}

			// Check if user in course. If not add
			var courseUser models.CourseUser
			{
				if err := transaction.
					Where("user_id = ?", user.ID).
					Where("course_id = ?", params.CourseID).
					Find(&courseUser).Error; err != nil {
					transaction.Rollback()
					return &common.ErrorResponse{
						Code:    500,
						Message: "Failed to remove student from class",
					}
				}

				if courseUser.ID == 0 {
					courseUser.CourseID = params.CourseID
					courseUser.UserID = user.ID
				}

				hasStudentRole := false
				for _, rl := range courseUser.Roles {
					if rl == enums.CourseUserRoleStudent {
						hasStudentRole = true
					}
				}
				if !hasStudentRole {
					courseUser.Roles = append(courseUser.Roles, enums.CourseUserRoleStudent)
				}

				if activityStudent.StudyFormCode == "P" {
					a := enums.StudyFormFulltime
					courseUser.StudyForm = &a
				} else {
					a := enums.StudyFormCombined
					courseUser.StudyForm = &a
				}

				if err := transaction.
					Save(&courseUser).Error; err != nil {
					transaction.Rollback()
					return &common.ErrorResponse{
						Code:    500,
						Message: "Failed to add user to course",
					}
				}
			}

			// Check if user in class. If already is. Return error. If not add him to this one
			var classUsers []models.ClassStudent
			{

				if err := transaction.
					Where("user_id = ?", user.ID).
					Where("class_id in ?", courseClassesIDs).
					Find(&classUsers).Error; err != nil {
					transaction.Rollback()
					return &common.ErrorResponse{
						Code:    500,
						Message: "Failed to remove student from class",
					}
				}

				shouldSkip := false
				for _, classUser := range classUsers {
					if classUser.ClassID != class.ID {
						if class.Type != enums.ClassTypeP {
							shouldSkip = true
							break
						}
					}
				}

				if shouldSkip {
					classAssignSkipped = append(classAssignSkipped, courseUser)
					fmt.Println("SKIPPING")
					continue
				}

				var classUser *models.ClassStudent
				if err := transaction.
					Where("user_id = ?", user.ID).
					Where("class_id = ?", class.ID).
					Find(&classUser).Error; err != nil {
					transaction.Rollback()
					return &common.ErrorResponse{
						Code:    500,
						Message: "Failed to remove student from class",
					}
				}

				if classUser.ID == 0 {
					classUser.UserID = user.ID
					classUser.ClassID = class.ID
				}

				if err := transaction.
					Save(&classUser).Error; err != nil {
					transaction.Rollback()
					return &common.ErrorResponse{
						Code:    500,
						Message: "Failed to add user to course",
					}
				}
			}
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to remove student from class",
		}
	}

	var classStudents []models.ClassStudent
	if err := initializers.DB.
		Where("class_id = ?", params.ClassID).
		InnerJoins("User").
		Find(&classStudents).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch class students",
			Details: err.Error(),
		}
	}

	utils.DebugPrintJSON(classAssignSkipped)

	// Convert to DTOs
	dtoList := make([]dtos.ClassUserDTO, len(classStudents))
	for i, c := range classStudents {
		dtoList[i] = dtos.ClassUserDTO{}.From(c.User)
	}

	c.JSON(200, RemoveStudentResponse{
		Students: dtoList,
	})
	return nil
}
