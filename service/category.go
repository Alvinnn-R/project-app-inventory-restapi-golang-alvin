package service

import (
	"errors"
	"project-app-inventory/dto"
	"project-app-inventory/model"
	"project-app-inventory/repository"
	"project-app-inventory/utils"
)

type CategoryService interface {
	Create(category *model.Category) error
	GetAllCategories(page, limit int) (*[]model.Category, *dto.Pagination, error)
	GetCategoryByID(id int) (*model.Category, error)
	Update(id int, data *model.Category) error
	Delete(id int) error
}

type categoryService struct {
	Repo repository.Repository
}

func NewCategoryService(repo repository.Repository) CategoryService {
	return &categoryService{Repo: repo}
}

func (s *categoryService) Create(category *model.Category) error {
	// Check if name already exists
	existingCategory, err := s.Repo.CategoryRepo.FindByName(category.Name)
	if err != nil {
		return errors.New("failed to check category name")
	}
	if existingCategory != nil {
		return errors.New("category name already exists")
	}

	return s.Repo.CategoryRepo.Create(category)
}

func (s *categoryService) GetAllCategories(page, limit int) (*[]model.Category, *dto.Pagination, error) {
	categories, total, err := s.Repo.CategoryRepo.FindAll(page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		CurrentPage:  page,
		Limit:        limit,
		TotalPages:   utils.TotalPage(limit, int64(total)),
		TotalRecords: total,
	}
	return &categories, &pagination, nil
}

func (s *categoryService) GetCategoryByID(id int) (*model.Category, error) {
	category, err := s.Repo.CategoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

func (s *categoryService) Update(id int, data *model.Category) error {
	// Check if category exists
	existingCategory, err := s.Repo.CategoryRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existingCategory == nil {
		return errors.New("category not found")
	}

	// If name is empty, keep existing name
	if data.Name == "" {
		data.Name = existingCategory.Name
	}

	// If description is nil, keep existing description
	if data.Description == nil {
		data.Description = existingCategory.Description
	}

	// Check if name is being changed and if new name already exists
	if data.Name != existingCategory.Name {
		nameExists, err := s.Repo.CategoryRepo.FindByName(data.Name)
		if err != nil {
			return errors.New("failed to check category name")
		}
		if nameExists != nil {
			return errors.New("category name already exists")
		}
	}

	return s.Repo.CategoryRepo.Update(id, data)
}

func (s *categoryService) Delete(id int) error {
	// Check if category exists
	existingCategory, err := s.Repo.CategoryRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existingCategory == nil {
		return errors.New("category not found")
	}

	return s.Repo.CategoryRepo.Delete(id)
}
