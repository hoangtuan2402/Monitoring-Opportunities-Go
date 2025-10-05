//go:build wireinject
// +build wireinject

package di

import (
	http "Monitoring-Opportunities/src/api"
	"Monitoring-Opportunities/src/api/controller"
	"Monitoring-Opportunities/src/config"
	"Monitoring-Opportunities/src/database"
	"Monitoring-Opportunities/src/gateway"
	"Monitoring-Opportunities/src/repository"
	service "Monitoring-Opportunities/src/services"

	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(
		// Database
		database.NewMongoDatabase,

		// Gateways
		gateway.NewEthereumRpcClient,
		gateway.NewUniswapV2Gateway,

		// Repositories
		repository.NewProductRepository,

		// Services
		service.NewUserService,
		service.NewProductService,
		service.NewPoolService,

		// Controllers
		handler.NewUserController,
		handler.NewProductController,
		handler.NewPoolController,

		// Server
		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil
}
