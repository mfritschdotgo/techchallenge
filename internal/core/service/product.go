package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mfritschdotgo/techchallenge/internal/adapter/handler/dto"
	"github.com/mfritschdotgo/techchallenge/internal/core/domain"
	"github.com/mfritschdotgo/techchallenge/internal/core/port"
)

type ProductService struct {
	productRepo     port.ProductRepository
	categoryService *CategoryService
}

func NewProductService(repo port.ProductRepository, categoryService *CategoryService) *ProductService {
	return &ProductService{
		productRepo:     repo,
		categoryService: categoryService,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, dto dto.CreateProductRequest) (*domain.Product, error) {

	if _, err := s.categoryService.GetCategoryByID(ctx, dto.CategoryId.String()); err != nil {
		return nil, fmt.Errorf("category validation failed: %w", err)
	}

	product, err := domain.NewProduct(dto.Name, dto.Price, dto.CategoryId, dto.Description, dto.Image)
	if err != nil {
		return nil, err
	}

	if _, err = s.productRepo.CreateProduct(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) ReplaceProduct(ctx context.Context, id string, productDto dto.CreateProductRequest) (*domain.Product, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	if _, err := s.categoryService.GetCategoryByID(ctx, productDto.CategoryId.String()); err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	product, err := s.productRepo.GetProductByID(ctx, uuidID.String())
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	product.Name = productDto.Name
	product.Price = productDto.Price
	product.CategoryId = productDto.CategoryId
	product.Description = productDto.Description
	product.Image = productDto.Image
	product.UpdatedAt = time.Now()

	if _, err := s.productRepo.ReplaceProduct(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id string, productDto dto.CreateProductRequest) (*domain.Product, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	product, err := s.productRepo.GetProductByID(ctx, uuidID.String())
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	if productDto.Name != "" {
		product.Name = productDto.Name
	}
	if productDto.Price != 0 {
		product.Price = productDto.Price
	}

	if productDto.CategoryId != uuid.Nil {
		if _, err := s.categoryService.GetCategoryByID(ctx, productDto.CategoryId.String()); err != nil {
			return nil, fmt.Errorf(err.Error())
		}

		product.CategoryId = productDto.CategoryId
	}
	if productDto.Description != "" {
		product.Description = productDto.Description
	}
	if productDto.Image != "" {
		product.Image = productDto.Image
	}

	product.UpdatedAt = time.Now()

	if _, err := s.productRepo.ReplaceProduct(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	product, err := s.productRepo.GetProductByID(ctx, uuidID.String())
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return product, nil
}

func (s *ProductService) GetProducts(ctx context.Context, category string, page, size int) ([]domain.Product, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	if category != "" {
		if _, err := s.categoryService.GetCategoryByID(ctx, category); err != nil {
			return nil, fmt.Errorf(err.Error())
		}
	}
	products, err := s.productRepo.GetProducts(ctx, category, page, size)
	if err != nil {
		return nil, fmt.Errorf("error retrieving products: %w", err)
	}

	return products, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	if err := s.productRepo.DeleteProduct(ctx, uuidID.String()); err != nil {
		return fmt.Errorf("product not found or error deleting product: %w", err)
	}

	return nil
}
