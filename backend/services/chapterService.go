package services

import (
	"elogika.vsb.cz/backend/repositories"
)

type ChapterService struct {
	courseRepo *repositories.CategoryRepository
}

func NewChapterService(repo *repositories.CategoryRepository) *ChapterService {
	return &ChapterService{courseRepo: repo}
}
