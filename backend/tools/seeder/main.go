package main

import (
	"fmt"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/tools/seeder/seed"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB(false)
}

func main() {
	users := seed.CreateUsers()
	courses := seed.CreateCourses()
	seed.CreateCourseUsers(courses, users)
	seed.CreateChapters(courses)

	fmt.Println("Everything created successfully")
}
