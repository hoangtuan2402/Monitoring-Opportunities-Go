package http

import (
	"Monitoring-Opportunities/src/api/controller"
	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userController *handler.UserController) *ServerHTTP {
	engine := gin.New()

	// Use logger from Gin
	engine.Use(gin.Logger())

	// Request JWT
	//engine.POST("/login", middleware.LoginHandler)

	// Auth middleware
	api := engine.Group("/api")

	api.GET("users", userController.GetByEmail)
	api.GET("users/:id", userController.GetByID)
	api.POST("users", userController.Create)
	api.DELETE("users/:id", userController.Delete)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run(":3000")
}
