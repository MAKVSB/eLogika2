package seed

import (
	"fmt"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

func CreateCourseUsers(courses []models.Course, users []models.User) []models.CourseUser {
	var course_users []models.CourseUser

	course_users = append(course_users, models.CourseUser{
		Roles:    []enums.CourseUserRoleEnum{enums.CourseUserRoleGarant},
		CourseID: courses[0].ID,
		UserID:   users[0].ID,
	})
	course_users = append(course_users, models.CourseUser{
		Roles:    []enums.CourseUserRoleEnum{enums.CourseUserRoleStudent},
		CourseID: courses[1].ID,
		UserID:   users[0].ID,
	})

	for index, course_user := range course_users {
		if err := initializers.DB.Create(&course_user).Error; err != nil {
			fmt.Println("Failed to insert", index, err)
		}
	}

	return course_users
}
