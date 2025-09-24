package services

import (
	"elogika.vsb.cz/backend/repositories"
)

type ClassService struct {
	courseRepo *repositories.ClassRepository
}

func NewClassService(repo *repositories.ClassRepository) *ClassService {
	return &ClassService{courseRepo: repo}
}
