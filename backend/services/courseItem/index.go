package services_course_item

import (
	"elogika.vsb.cz/backend/repositories"
)

type CourseItemService struct {
	courseItemRepo *repositories.CourseItemRepository
}

func NewCourseItemService(repo *repositories.CourseItemRepository) *CourseItemService {
	return &CourseItemService{courseItemRepo: repo}
}
