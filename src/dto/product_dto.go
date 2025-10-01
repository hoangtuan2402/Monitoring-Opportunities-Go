package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProductDTO struct {
	UUID        uuid.UUID `json:"id" example:"a53515e3-5a7f-440b-82f6-3d84ac7ce746"`
	Name        string    `json:"name" example:"iPhone 15 Pro"`
	Description string    `json:"description" example:"Latest iPhone model"`
	Price       float64   `json:"price" example:"999.99"`
	Stock       int       `json:"stock" example:"100"`
}

type CreateProduct struct {
	Name        string    `json:"name" binding:"required" example:"iPhone 15 Pro"`
	Description string    `json:"description" example:"Latest iPhone model"`
	Price       float64   `json:"price" binding:"required" example:"999.99"`
	Stock       int       `json:"stock" binding:"required" example:"100"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type UpdateProduct struct {
	Name        string  `json:"name" example:"iPhone 15 Pro"`
	Description string  `json:"description" example:"Latest iPhone model"`
	Price       float64 `json:"price" example:"999.99"`
	Stock       int     `json:"stock" example:"100"`
}
