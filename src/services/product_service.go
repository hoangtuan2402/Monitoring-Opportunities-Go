package service

import (
	"Monitoring-Opportunities/src/dto"
	"Monitoring-Opportunities/src/models"
	"Monitoring-Opportunities/src/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrProductNotFound       = errors.New("product not found")
	ErrCreateProductValidate = errors.New("failed to create a product")
	ErrUpdateProductValidate = errors.New("failed to update a product")
	ErrDeleteProduct         = errors.New("failed to delete a product")
	ErrDuplicateProductName  = errors.New("product with that name already exists")
)

type ProductService interface {
	GetAll() ([]dto.ProductDTO, error)
	Create(product dto.CreateProduct) (dto.ProductDTO, error)
	Update(product dto.UpdateProduct, productID uuid.UUID) (dto.ProductDTO, error)
	Delete(id uuid.UUID) (dto.ProductDTO, error)
	FindByID(id uuid.UUID) (dto.ProductDTO, error)
	FindByName(name string) ([]dto.ProductDTO, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) GetAll() ([]dto.ProductDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	products, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return s.toProductDTOList(products), nil
}

func (s *productService) Create(form dto.CreateProduct) (dto.ProductDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	product := &models.Product{
		Name:        form.Name,
		Description: form.Description,
		Price:       form.Price,
		Stock:       form.Stock,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return dto.ProductDTO{}, ErrCreateProductValidate
	}

	return s.toProductDTO(product), nil
}

func (s *productService) Update(updateData dto.UpdateProduct, productID uuid.UUID) (dto.ProductDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	existingProduct, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return dto.ProductDTO{}, ErrProductNotFound
		}
		return dto.ProductDTO{}, ErrUpdateProductValidate
	}

	existingProduct.Name = updateData.Name
	existingProduct.Description = updateData.Description
	existingProduct.Price = updateData.Price
	existingProduct.Stock = updateData.Stock

	if err := s.repo.Update(ctx, existingProduct); err != nil {
		return dto.ProductDTO{}, ErrUpdateProductValidate
	}

	return s.toProductDTO(existingProduct), nil
}

func (s *productService) Delete(id uuid.UUID) (dto.ProductDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return dto.ProductDTO{}, ErrProductNotFound
		}
		return dto.ProductDTO{}, ErrDeleteProduct
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return dto.ProductDTO{}, ErrDeleteProduct
	}

	return s.toProductDTO(product), nil
}

func (s *productService) FindByID(id uuid.UUID) (dto.ProductDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return dto.ProductDTO{}, ErrProductNotFound
		}
		return dto.ProductDTO{}, err
	}

	return s.toProductDTO(product), nil
}

func (s *productService) FindByName(name string) ([]dto.ProductDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	products, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return s.toProductDTOList(products), nil
}

// Helper functions to convert models to DTOs
func (s *productService) toProductDTO(product *models.Product) dto.ProductDTO {
	return dto.ProductDTO{
		UUID:        product.UUID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}
}

func (s *productService) toProductDTOList(products []models.Product) []dto.ProductDTO {
	dtos := make([]dto.ProductDTO, len(products))
	for i, product := range products {
		dtos[i] = dto.ProductDTO{
			UUID:        product.UUID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
		}
	}
	return dtos
}
