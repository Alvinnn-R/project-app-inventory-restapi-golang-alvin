package service

import (
	"errors"
	"project-app-inventory/dto"
	"project-app-inventory/model"
	"project-app-inventory/repository"
	"project-app-inventory/utils"
)

type WarehouseService interface {
	Create(warehouse *model.Warehouse) error
	GetAllWarehouses(page, limit int) (*[]model.Warehouse, *dto.Pagination, error)
	GetWarehouseByID(id int) (*model.Warehouse, error)
	Update(id int, data *model.Warehouse) error
	Delete(id int) error
}

type warehouseService struct {
	Repo repository.Repository
}

func NewWarehouseService(repo repository.Repository) WarehouseService {
	return &warehouseService{Repo: repo}
}

func (s *warehouseService) Create(warehouse *model.Warehouse) error {
	// Check if name already exists
	existingWarehouse, err := s.Repo.WarehouseRepo.FindByName(warehouse.Name)
	if err != nil {
		return errors.New("failed to check warehouse name")
	}
	if existingWarehouse != nil {
		return errors.New("warehouse name already exists")
	}

	return s.Repo.WarehouseRepo.Create(warehouse)
}

func (s *warehouseService) GetAllWarehouses(page, limit int) (*[]model.Warehouse, *dto.Pagination, error) {
	warehouses, total, err := s.Repo.WarehouseRepo.FindAll(page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		CurrentPage:  page,
		Limit:        limit,
		TotalPages:   utils.TotalPage(limit, int64(total)),
		TotalRecords: total,
	}
	return &warehouses, &pagination, nil
}

func (s *warehouseService) GetWarehouseByID(id int) (*model.Warehouse, error) {
	warehouse, err := s.Repo.WarehouseRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, errors.New("warehouse not found")
	}
	return warehouse, nil
}

func (s *warehouseService) Update(id int, data *model.Warehouse) error {
	// Check if warehouse exists
	existingWarehouse, err := s.Repo.WarehouseRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existingWarehouse == nil {
		return errors.New("warehouse not found")
	}

	// If name is empty, keep existing name
	if data.Name == "" {
		data.Name = existingWarehouse.Name
	}

	// If location is empty, keep existing location
	if data.Location == "" {
		data.Location = existingWarehouse.Location
	}

	// Check if name is being changed and if new name already exists
	if data.Name != existingWarehouse.Name {
		nameExists, err := s.Repo.WarehouseRepo.FindByName(data.Name)
		if err != nil {
			return errors.New("failed to check warehouse name")
		}
		if nameExists != nil {
			return errors.New("warehouse name already exists")
		}
	}

	return s.Repo.WarehouseRepo.Update(id, data)
}

func (s *warehouseService) Delete(id int) error {
	// Check if warehouse exists
	existingWarehouse, err := s.Repo.WarehouseRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existingWarehouse == nil {
		return errors.New("warehouse not found")
	}

	return s.Repo.WarehouseRepo.Delete(id)
}
