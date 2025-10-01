package handler

import (
	"Monitoring-Opportunities/src/common"
	"Monitoring-Opportunities/src/dto"
	service "Monitoring-Opportunities/src/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type ProductController struct {
	productService service.ProductService
}

func NewProductController(productService service.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

func (c *ProductController) GetAll(ctx *gin.Context) {
	products, err := c.productService.GetAll()
	if err != nil {
		log.Printf("Failed to get products: %v", err)

		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "failed to get products data"})
		return
	}

	ctx.JSON(
		http.StatusOK,
		common.BaseResponse[[]dto.ProductDTO]{
			Status:  http.StatusOK,
			Message: "Successfully get product data",
			Data:    products,
		},
	)
}

func (c *ProductController) Create(ctx *gin.Context) {
	var product dto.CreateProduct
	if err := ctx.ShouldBindBodyWithJSON(&product); err != nil {
		log.Printf("Failed to create product: %v", err)

		ctx.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%s: %s", service.ErrCreateProductValidate, err.Error())})
		return
	}

	createdProduct, _ := c.productService.Create(product)

	ctx.JSON(
		http.StatusCreated,
		common.BaseResponse[dto.ProductDTO]{
			Status:  http.StatusCreated,
			Message: "Product created successfully",
			Data:    createdProduct,
		},
	)
}

func (c *ProductController) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	parsedUUID, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("Failed to parse product ID: %v", err)

		ctx.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("invalid product ID: %s", err.Error())})
		return
	}

	var product dto.UpdateProduct
	if err := ctx.ShouldBindBodyWithJSON(&product); err != nil {
		log.Printf("Failed to update product: %v", err)

		ctx.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("failed to update a product: %s", err.Error())})
		return
	}

	updatedProduct, _ := c.productService.Update(product, parsedUUID)

	ctx.JSON(
		http.StatusOK,
		common.BaseResponse[dto.ProductDTO]{
			Status:  http.StatusOK,
			Message: "Product updated successfully",
			Data:    updatedProduct,
		},
	)
}

func (c *ProductController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	parsedUUID, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("Failed to parse product ID: %v", err)

		ctx.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("invalid product ID: %s", err.Error())})
		return
	}

	deletedProduct, err := c.productService.Delete(parsedUUID)
	if err != nil {
		log.Printf("Failed to delete product: %v", err)

		switch err {
		case service.ErrProductNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"errors": fmt.Sprintf("%s: %s", service.ErrProductNotFound, err.Error())})
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("failed to delete a product: %s", err.Error())})
		}

		return
	}

	ctx.JSON(
		http.StatusOK,
		common.BaseResponse[dto.ProductDTO]{
			Status:  http.StatusOK,
			Message: "Product deleted successfully",
			Data:    deletedProduct,
		},
	)
}

func (c *ProductController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	parsedUUID, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("Failed to parse product ID: %v", err)

		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	product, err := c.productService.FindByID(parsedUUID)
	if err != nil {
		log.Printf("Failed to find product with id %s: %v", parsedUUID, err)

		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	ctx.JSON(
		http.StatusOK,
		common.BaseResponse[dto.ProductDTO]{
			Status:  http.StatusOK,
			Message: fmt.Sprintf("Successfully get product with id %s", product.UUID),
			Data:    product,
		},
	)
}

func (c *ProductController) GetByName(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "name query parameter is required"})
		return
	}

	products, err := c.productService.FindByName(name)
	if err != nil {
		log.Printf("Failed to find products with name %s: %v", name, err)

		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	ctx.JSON(
		http.StatusOK,
		common.BaseResponse[[]dto.ProductDTO]{
			Status:  http.StatusOK,
			Message: fmt.Sprintf("Successfully get products with name %s", name),
			Data:    products,
		},
	)
}
