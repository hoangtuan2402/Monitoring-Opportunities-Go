package service

import (
	"Monitoring-Opportunities/src/dto"
	"errors"
	"github.com/google/uuid"
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
}

func NewProductService() ProductService {
	return &productService{}
}

func (s *productService) GetAll() ([]dto.ProductDTO, error) {
	// Mock data
	return []dto.ProductDTO{
		{
			UUID:        uuid.New(),
			Name:        "iPhone 15 Pro",
			Description: "Latest iPhone model",
			Price:       999.99,
			Stock:       100,
		},
		{
			UUID:        uuid.New(),
			Name:        "Samsung Galaxy S24",
			Description: "Latest Samsung flagship",
			Price:       899.99,
			Stock:       50,
		},
	}, nil
}

func (s *productService) Create(form dto.CreateProduct) (dto.ProductDTO, error) {
	return dto.ProductDTO{
		UUID:        uuid.New(),
		Name:        form.Name,
		Description: form.Description,
		Price:       form.Price,
		Stock:       form.Stock,
	}, nil
}

func (s *productService) Update(product dto.UpdateProduct, productID uuid.UUID) (dto.ProductDTO, error) {
	return dto.ProductDTO{
		UUID:        productID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}, nil
}

func (s *productService) Delete(id uuid.UUID) (dto.ProductDTO, error) {
	return dto.ProductDTO{
		UUID:        id,
		Name:        "Deleted Product",
		Description: "Product has been deleted",
		Price:       0,
		Stock:       0,
	}, nil
}

func (s *productService) FindByID(id uuid.UUID) (dto.ProductDTO, error) {
	return dto.ProductDTO{
		UUID:        id,
		Name:        "Product Found By ID",
		Description: "Product description",
		Price:       499.99,
		Stock:       75,
	}, nil
}

func (s *productService) FindByName(name string) ([]dto.ProductDTO, error) {
	return []dto.ProductDTO{
		{
			UUID:        uuid.New(),
			Name:        name,
			Description: "Product found by name",
			Price:       599.99,
			Stock:       30,
		},
	}, nil
}
