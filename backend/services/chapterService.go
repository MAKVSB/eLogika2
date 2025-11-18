package services

import (
	"elogika.vsb.cz/backend/repositories"
)

type ChapterService struct {
	courseRepo *repositories.ChapterRepository
}

func NewChapterService(repo *repositories.ChapterRepository) *ChapterService {
	return &ChapterService{courseRepo: repo}
}
