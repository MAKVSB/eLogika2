package services

import "elogika.vsb.cz/backend/repositories"

type CategoryService struct {
	courseRepo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{courseRepo: repo}
}
