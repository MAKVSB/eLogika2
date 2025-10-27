package handlers

import (
	"fmt"
	"strconv"
	"strings"

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
	"gorm.io/gorm"
)

// @Description Updated list of tutors
type ClassImportStudentsResponse struct {
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

	classService := services.NewClassService(repositories.NewClassRepository())
	class, err := classService.GetClassByID(initializers.DB, params.CourseID, params.ClassID, userData.ID, userRole, nil, false, nil)
	if err != nil {
		return err
	}

	transaction := initializers.DB.Begin()

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

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to remove student from class",
		}
	}

	c.JSON(200, ClassImportStudentsResponse{
		Success: true,
	})
	return nil
}

func ImportByConcreteActivity(dbRef *gorm.DB, course *models.Course, userRole enums.CourseUserRoleEnum, class *models.Class) *common.ErrorResponse {
	var courseClasses []*models.Class
	if err := dbRef.
		Where("course_id = ?", course.ID).
		Where("type != ?", enums.ClassTypeP).
		Where("id != ?", class.ID).
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
	if concreteActivities == nil || len(*concreteActivities) == 0 {
		return &common.ErrorResponse{
			Code:    0,
			Message: "Failed load information from inbus",
			Details: "Schedule is missing",
		}
	}
	concreteActivity, err := GetMatchingActivity(*concreteActivities, class)
	if err != nil {
		return err
	}

	concreteActivityStudents, err := inbusClient.GetConcreteActivityStudents(concreteActivity.ConcreteActivityId)
	if err != nil {
		dbRef.Rollback()
		return &common.ErrorResponse{
			Code:    409,
			Message: "Failed to get students assignned to schedule item",
		}
	}

	for _, activityStudent := range *concreteActivityStudents {
		// If user does in exist in entire system, create. If does, update
		user, err := UpsertUser(dbRef, activityStudent)
		if err != nil {
			return err
		}

		// Check if user in course. If not add
		_, err = UpsertCourseUser(dbRef, course.ID, user.ID, activityStudent.StudyFormCode)
		if err != nil {
			return err
		}

		// Check if user in class. If already is. Return error. If not add him to this one
		_, err = UpsertClassUser(dbRef, user, courseClassesIDs, class)
		if err != nil {
			return err
		}
	}
	return nil
}

func ImportPartTimeStudents(dbRef *gorm.DB, course *models.Course, userRole enums.CourseUserRoleEnum, class *models.Class) *common.ErrorResponse {
	var courseClasses []*models.Class
	if err := dbRef.
		Where("course_id = ?", course.ID).
		Where("type != ?", enums.ClassTypeP).
		Where("id != ?", class.ID).
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

	subjectVersionStudents, err := inbusClient.GetSubjectVersionStudents(subject.SubjectVersionId, semester.SemesterId)
	if err != nil {
		return err
	}

	for _, courseStudent := range *subjectVersionStudents {
		if courseStudent.StudyFormCode != "K" {
			continue
		}
		// If user does in exist in entire system, create. If does, update
		user, err := UpsertUser(dbRef, courseStudent)
		if err != nil {
			return err
		}

		// Check if user in course. If not add
		_, err = UpsertCourseUser(dbRef, course.ID, user.ID, courseStudent.StudyFormCode)
		if err != nil {
			return err
		}

		// Check if user in class. If already is. Return error. If not add him to this one
		_, err = UpsertClassUser(dbRef, user, courseClassesIDs, class)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpsertUser(dbRef *gorm.DB, activityStudent *inbus.StudyRelation) (*models.User, *common.ErrorResponse) {
	// If user does in exist in entire system, create. If does, update
	var user *models.User
	if err := dbRef.
		Where("identity_provider_id = ?", activityStudent.PersonId).
		Find(&user).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to remove student from class",
		}
	}

	user.Version = user.Version + 1
	user.Username = activityStudent.Login
	user.DegreeBefore = activityStudent.DegreeBefore
	user.FirstName = activityStudent.FirstName
	user.FamilyName = activityStudent.SecondName
	user.DegreeAfter = activityStudent.DegreeAfter
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

	if err := dbRef.
		Save(&user).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to save new user",
		}
	}

	return user, nil
}

func UpsertCourseUser(dbRef *gorm.DB, courseId uint, userId uint, studyFormCode string) (*models.CourseUser, *common.ErrorResponse) {
	var courseUser *models.CourseUser
	if err := dbRef.
		Where("user_id = ?", userId).
		Where("course_id = ?", courseId).
		Find(&courseUser).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to remove student from class",
		}
	}

	if courseUser.ID == 0 {
		courseUser.CourseID = courseId
		courseUser.UserID = userId
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

	if studyFormCode == "P" {
		a := enums.StudyFormFulltime
		courseUser.StudyForm = &a
	} else {
		a := enums.StudyFormCombined
		courseUser.StudyForm = &a
	}

	if err := dbRef.
		Save(&courseUser).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to add user to course",
		}
	}
	return courseUser, nil
}

func UpsertClassUser(dbRef *gorm.DB, user *models.User, courseClassIDs []uint, currentClass *models.Class) (*models.ClassStudent, *common.ErrorResponse) {
	var userClasses []models.ClassStudent

	// If class type is "přednáška" user can be imported regardless
	if currentClass.Type != enums.ClassTypeP {
		if err := dbRef.
			Where("user_id = ?", user.ID).
			Where("class_id in ?", courseClassIDs).
			Where("class_id != ?", currentClass.ID).
			Find(&userClasses).Error; err != nil {
			return nil, &common.ErrorResponse{
				Code:    500,
				Message: "Failed to find class students",
			}
		}
	}

	if len(userClasses) != 0 {
		res := []common.ErrorResources{}

		for _, uc := range userClasses {
			res = append(res, common.ErrorResources{
				ResourceType: "class",
				ResourceID:   uc.ClassID,
			})
		}

		res = append(res, common.ErrorResources{
			ResourceType: "user",
			ResourceData: user.FullName(),
		})

		return nil, &common.ErrorResponse{
			Code:      409,
			Message:   "Failed to assign user to class",
			Details:   "User is already member of another class" + strconv.Itoa(int(currentClass.ID)),
			Resources: res,
		}
	}

	var classUser *models.ClassStudent
	if err := dbRef.
		Where("user_id = ?", user.ID).
		Where("class_id = ?", currentClass.ID).
		Find(&classUser).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to remove student from class",
		}
	}

	if classUser.ID == 0 {
		classUser.UserID = user.ID
		classUser.ClassID = currentClass.ID
	}

	if err := dbRef.
		Save(&classUser).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to add user to course",
		}
	}

	return classUser, nil
}

func GetMatchingActivity(concreteActivities []*inbus.ConcreteActivity, class *models.Class) (*inbus.ConcreteActivity, *common.ErrorResponse) {
	var matchingActivity *inbus.ConcreteActivity
	for _, concreteActivity := range concreteActivities {
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

		if !strings.HasPrefix(class.Room, "POR") {
			class.Room = "POR" + class.Room
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

		if matchingActivity == nil {
			matchingActivity = concreteActivity
		} else {
			return nil, &common.ErrorResponse{
				Code:    409,
				Message: "Multiple schedule items matched",
			}
		}

		if class.StudyForm == enums.StudyFormFulltime {
			// Class type
			classTypeId := "P"
			switch class.Type {
			case enums.ClassTypeC:
				classTypeId = "C"
			case enums.ClassTypeP:
				classTypeId = "P"
			default:
				panic(fmt.Sprintf("unexpected enums.ClassTypeEnum: %#v", class.Type))
			}
			if classTypeId != concreteActivity.EducationTypeAbbrev {
				continue
			}
		}
	}

	if matchingActivity == nil {
		return nil, &common.ErrorResponse{
			Code:    409,
			Message: "No schedule item matched",
		}
	} else {
		return matchingActivity, nil
	}
}
