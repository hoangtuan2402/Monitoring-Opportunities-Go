//go:build wireinject
// +build wireinject

package di

import (
	http "Monitoring-Opportunities/src/api"
	"Monitoring-Opportunities/src/api/controller"
	"Monitoring-Opportunities/src/config"
	service "Monitoring-Opportunities/src/services"
	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(
		service.NewUserService,
		handler.NewUserController,
		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil
}
