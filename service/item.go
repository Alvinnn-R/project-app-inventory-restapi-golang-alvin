package service

import (
	"errors"
	"project-app-inventory/dto"
	"project-app-inventory/model"
	"project-app-inventory/repository"
	"project-app-inventory/utils"
)

type ItemService interface {
	Create(item *model.Item) error
	GetAllItems(page, limit int) (*[]model.Item, *dto.Pagination, error)
	GetItemByID(id int) (*model.Item, error)
	Update(id int, data *model.Item) error
	Delete(id int) error
}

type itemService struct {
	Repo repository.Repository
}

func NewItemService(repo repository.Repository) ItemService {
	return &itemService{Repo: repo}
}

func (s *itemService) Create(item *model.Item) error {
	// Check if SKU already exists
	existingItem, err := s.Repo.ItemRepo.FindBySKU(item.SKU)
	if err != nil {
		return errors.New("failed to check SKU")
	}
	if existingItem != nil {
		return errors.New("SKU already exists")
	}

	return s.Repo.ItemRepo.Create(item)
}

func (s *itemService) GetAllItems(page, limit int) (*[]model.Item, *dto.Pagination, error) {
	items, total, err := s.Repo.ItemRepo.FindAll(page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		CurrentPage:  page,
		Limit:        limit,
		TotalPages:   utils.TotalPage(limit, int64(total)),
		TotalRecords: total,
	}
	return &items, &pagination, nil
}

func (s *itemService) GetItemByID(id int) (*model.Item, error) {
	item, err := s.Repo.ItemRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("item not found")
	}
	return item, nil
}

func (s *itemService) Update(id int, data *model.Item) error {
	// Check if item exists
	existingItem, err := s.Repo.ItemRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existingItem == nil {
		return errors.New("item not found")
	}

	// If fields are empty/zero, keep existing values
	if data.SKU == "" {
		data.SKU = existingItem.SKU
	}
	if data.Name == "" {
		data.Name = existingItem.Name
	}
	if data.CategoryID == 0 {
		data.CategoryID = existingItem.CategoryID
	}
	if data.RackID == 0 {
		data.RackID = existingItem.RackID
	}
	// Note: Stock, MinimumStock, and Price can be 0, so we don't check for zero values
	// If you want to keep existing values when 0 is sent, uncomment below:
	// if data.Stock == 0 {
	// 	data.Stock = existingItem.Stock
	// }
	// if data.MinimumStock == 0 {
	// 	data.MinimumStock = existingItem.MinimumStock
	// }
	// if data.Price == 0 {
	// 	data.Price = existingItem.Price
	// }

	// Check if SKU is being changed and if new SKU already exists
	if data.SKU != existingItem.SKU {
		skuExists, err := s.Repo.ItemRepo.FindBySKU(data.SKU)
		if err != nil {
			return errors.New("failed to check SKU")
		}
		if skuExists != nil {
			return errors.New("SKU already exists")
		}
	}

	return s.Repo.ItemRepo.Update(id, data)
}

func (s *itemService) Delete(id int) error {
	// Check if item exists
	existingItem, err := s.Repo.ItemRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existingItem == nil {
		return errors.New("item not found")
	}

	return s.Repo.ItemRepo.Delete(id)
}
