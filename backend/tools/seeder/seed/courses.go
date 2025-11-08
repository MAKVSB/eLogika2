package seed

import (
	"fmt"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/tools/seeder/consts"
)

func CreateCourses() []models.Course {
	var courses []models.Course

	courses = append(courses, models.Course{
		Version:   1,
		Name:      "Matematika pro zpracování znalostí",
		Shortname: "MPZZ",
		Public:    true,
		Year:      2025,
		Semester:  "SUMMER",
		Content:   consts.DefaultContent,
	})

	courses = append(courses, models.Course{
		Version:   1,
		Name:      "Úvod do logického myšlení",
		Shortname: "ULM",
		Public:    true,
		Year:      2025,
		Semester:  "WINTER",
		Content:   consts.DefaultContent,
	})

	for index, course := range courses {
		if err := initializers.DB.Create(&course).Error; err != nil {
			fmt.Println("Failed to insert", index, err)
		}
	}

	return courses
}
