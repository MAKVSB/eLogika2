package seed

import (
	"fmt"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/tools/seeder/consts"
)

func CreateChapters(courses []models.Course) []models.Chapter {
	var chapters []models.Chapter

	chapters = append(chapters, models.Chapter{
		Version:  1,
		Name:     "Matematika pro zpracování znalostí",
		Content:  consts.DefaultContent,
		CourseID: courses[0].ID,
		Visible:  false,
		Order:    0,
	})

	chapters = append(chapters, models.Chapter{
		Version:  1,
		Name:     "Úvod do logického myšlení",
		Content:  consts.DefaultContent,
		CourseID: courses[1].ID,
		Visible:  false,
		Order:    0,
	})

	for index, chapter := range chapters {
		if err := initializers.DB.Create(&chapter).Error; err != nil {
			fmt.Println("Failed to insert", index, err)
		}
	}

	return chapters
}
