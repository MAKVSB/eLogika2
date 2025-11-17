package main

import (
	"fmt"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB(true)
}

func main() {
	err := initializers.DB.AutoMigrate(
		&models.Question{},
		&models.Answer{},
		&models.QuestionCheck{},
		&models.CourseQuestion{},
		&models.QuestionAnswer{},
		&models.User{},
		&models.Course{},
		&models.CourseUser{},
		&models.AuthToken{},
		&models.File{},
		&models.Step{},
		&models.Chapter{},
		&models.Category{},

		&models.CourseItemGroup{},
		&models.CourseItemActivity{},
		&models.CourseItemTest{},
		&models.CourseItem{},

		&models.TemplateBlockSegment{},
		&models.TemplateBlock{},
		&models.Template{},
		&models.Term{},
		&models.UserTerm{},

		&models.Test{},
		&models.TestQuestion{},
		&models.TestQuestionAnswer{},

		&models.TestInstance{},
		&models.TestInstanceQuestion{},
		&models.TestInstanceQuestionAnswer{},

		&models.CourseItemResult{},
		&models.Class{},
		&models.ClassStudent{},
		&models.ClassTutor{},
		&models.TestInstanceEvent{},
		&models.ActivityInstance{},
		&models.Email{},
		&models.RecognizerFile{},
		&models.LogAccess{},
		&models.LogError{},
	)
	fmt.Println(err)
}
