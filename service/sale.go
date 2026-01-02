package service

import (
	"errors"
	"project-app-inventory/dto"
	"project-app-inventory/model"
	"project-app-inventory/repository"
	"project-app-inventory/utils"
)

type SaleService interface {
	Create(userID int, items []dto.SaleItemRequest) (*model.Sale, error)
	GetAllSales(page, limit int) (*[]model.Sale, *dto.Pagination, error)
	GetSaleByID(id int) (*model.Sale, []model.SaleItem, error)
	Update(id int, items []dto.SaleItemRequest) error
	Delete(id int) error
}

type saleService struct {
	Repo repository.Repository
}

func NewSaleService(repo repository.Repository) SaleService {
	return &saleService{Repo: repo}
}

func (s *saleService) Create(userID int, items []dto.SaleItemRequest) (*model.Sale, error) {
	if len(items) == 0 {
		return nil, errors.New("sale must have at least one item")
	}

	// Prepare sale items and calculate total
	var saleItems []model.SaleItem
	var totalAmount float64

	for _, item := range items {
		// Get item details
		itemData, err := s.Repo.ItemRepo.FindByID(item.ItemID)
		if err != nil {
			return nil, err
		}
		if itemData == nil {
			return nil, errors.New("item not found")
		}

		// Check stock availability
		if itemData.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for item: " + itemData.Name)
		}

		subtotal := itemData.Price * float64(item.Quantity)
		totalAmount += subtotal

		saleItems = append(saleItems, model.SaleItem{
			ItemID:      item.ItemID,
			Quantity:    item.Quantity,
			PriceAtSale: itemData.Price,
			Subtotal:    subtotal,
		})
	}

	// Create sale
	sale := &model.Sale{
		UserID:      userID,
		TotalAmount: totalAmount,
	}

	err := s.Repo.SaleRepo.Create(sale, saleItems)
	if err != nil {
		return nil, err
	}

	return sale, nil
}

func (s *saleService) GetAllSales(page, limit int) (*[]model.Sale, *dto.Pagination, error) {
	sales, total, err := s.Repo.SaleRepo.FindAll(page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		CurrentPage:  page,
		Limit:        limit,
		TotalPages:   utils.TotalPage(limit, int64(total)),
		TotalRecords: total,
	}
	return &sales, &pagination, nil
}

func (s *saleService) GetSaleByID(id int) (*model.Sale, []model.SaleItem, error) {
	sale, err := s.Repo.SaleRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	if sale == nil {
		return nil, nil, errors.New("sale not found")
	}

	items, err := s.Repo.SaleRepo.FindSaleItems(id)
	if err != nil {
		return nil, nil, err
	}

	return sale, items, nil
}

func (s *saleService) Update(id int, items []dto.SaleItemRequest) error {
	if len(items) == 0 {
		return errors.New("sale must have at least one item")
	}

	// Check if sale exists
	existingSale, err := s.Repo.SaleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existingSale == nil {
		return errors.New("sale not found")
	}

	// Prepare sale items and calculate total
	var saleItems []model.SaleItem
	var totalAmount float64

	for _, item := range items {
		// Get item details
		itemData, err := s.Repo.ItemRepo.FindByID(item.ItemID)
		if err != nil {
			return err
		}
		if itemData == nil {
			return errors.New("item not found")
		}

		subtotal := itemData.Price * float64(item.Quantity)
		totalAmount += subtotal

		saleItems = append(saleItems, model.SaleItem{
			ItemID:      item.ItemID,
			Quantity:    item.Quantity,
			PriceAtSale: itemData.Price,
			Subtotal:    subtotal,
		})
	}

	// Update sale
	sale := &model.Sale{
		ID:          id,
		TotalAmount: totalAmount,
	}

	err = s.Repo.SaleRepo.Update(id, sale, saleItems)
	if err != nil {
		return err
	}

	return nil
}

func (s *saleService) Delete(id int) error {
	// Check if sale exists
	existingSale, err := s.Repo.SaleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existingSale == nil {
		return errors.New("sale not found")
	}

	return s.Repo.SaleRepo.Delete(id)
}
