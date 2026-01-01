package service

import (
	"errors"
	"project-app-inventory/dto"
	"project-app-inventory/model"
	"project-app-inventory/repository"
	"project-app-inventory/utils"
)

type RackService interface {
	Create(rack *model.Rack) error
	GetAllRacks(page, limit int) (*[]model.Rack, *dto.Pagination, error)
	GetRacksByWarehouse(warehouseID, page, limit int) (*[]model.Rack, *dto.Pagination, error)
	GetRackByID(id int) (*model.Rack, error)
	Update(id int, data *model.Rack) error
	Delete(id int) error
}

type rackService struct {
	Repo repository.Repository
}

func NewRackService(repo repository.Repository) RackService {
	return &rackService{Repo: repo}
}

func (s *rackService) Create(rack *model.Rack) error {
	// Check if rack code already exists in the same warehouse
	existingRack, err := s.Repo.RackRepo.FindByWarehouseAndCode(rack.WarehouseID, rack.Code)
	if err != nil {
		return errors.New("failed to check rack code")
	}
	if existingRack != nil {
		return errors.New("rack code already exists in this warehouse")
	}

	return s.Repo.RackRepo.Create(rack)
}

func (s *rackService) GetAllRacks(page, limit int) (*[]model.Rack, *dto.Pagination, error) {
	racks, total, err := s.Repo.RackRepo.FindAll(page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		CurrentPage:  page,
		Limit:        limit,
		TotalPages:   utils.TotalPage(limit, int64(total)),
		TotalRecords: total,
	}
	return &racks, &pagination, nil
}

func (s *rackService) GetRacksByWarehouse(warehouseID, page, limit int) (*[]model.Rack, *dto.Pagination, error) {
	racks, total, err := s.Repo.RackRepo.FindByWarehouseID(warehouseID, page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		CurrentPage:  page,
		Limit:        limit,
		TotalPages:   utils.TotalPage(limit, int64(total)),
		TotalRecords: total,
	}
	return &racks, &pagination, nil
}

func (s *rackService) GetRackByID(id int) (*model.Rack, error) {
	rack, err := s.Repo.RackRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if rack == nil {
		return nil, errors.New("rack not found")
	}
	return rack, nil
}

func (s *rackService) Update(id int, data *model.Rack) error {
	// Check if rack exists
	existingRack, err := s.Repo.RackRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existingRack == nil {
		return errors.New("rack not found")
	}

	// If warehouse_id is 0, keep existing warehouse_id
	if data.WarehouseID == 0 {
		data.WarehouseID = existingRack.WarehouseID
	}

	// If code is empty, keep existing code
	if data.Code == "" {
		data.Code = existingRack.Code
	}

	// If description is nil, keep existing description
	if data.Description == nil {
		data.Description = existingRack.Description
	}

	// Check if code is being changed or warehouse changed and if new combination already exists
	if data.Code != existingRack.Code || data.WarehouseID != existingRack.WarehouseID {
		codeExists, err := s.Repo.RackRepo.FindByWarehouseAndCode(data.WarehouseID, data.Code)
		if err != nil {
			return errors.New("failed to check rack code")
		}
		if codeExists != nil && codeExists.ID != id {
			return errors.New("rack code already exists in this warehouse")
		}
	}

	return s.Repo.RackRepo.Update(id, data)
}

func (s *rackService) Delete(id int) error {
	// Check if rack exists
	existingRack, err := s.Repo.RackRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existingRack == nil {
		return errors.New("rack not found")
	}

	return s.Repo.RackRepo.Delete(id)
}
