package http

import (
	"Monitoring-Opportunities/src/api/controller"
	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(
	userController *handler.UserController,
	productController *handler.ProductController,
) *ServerHTTP {
	engine := gin.New()

	// Use logger from Gin
	engine.Use(gin.Logger())

	// Request JWT
	//engine.POST("/login", middleware.LoginHandler)

	// Auth middleware
	api := engine.Group("/api")

	// User routes
	api.GET("users", userController.GetByEmail)
	api.GET("users/:id", userController.GetByID)
	api.POST("users", userController.Create)
	api.DELETE("users/:id", userController.Delete)

	// Product routes
	api.GET("products", productController.GetAll)
	api.GET("products/search", productController.GetByName)
	api.GET("products/:id", productController.GetByID)
	api.POST("products", productController.Create)
	api.PUT("products/:id", productController.Update)
	api.DELETE("products/:id", productController.Delete)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run(":3000")
}
