package service

import (
	"Monitoring-Opportunities/src/gateway"
	"Monitoring-Opportunities/src/models"
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidPoolAddress = errors.New("invalid pool address")
	ErrFetchPoolData      = errors.New("failed to fetch pool data")
)

type PoolService interface {
	GetPoolData(poolAddress string) (*models.PoolData, error)
}

type poolService struct {
	dexGateway gateway.DEXGateway
}

func NewPoolService(dexGateway gateway.DEXGateway) PoolService {
	return &poolService{
		dexGateway: dexGateway,
	}
}

func (s *poolService) GetPoolData(poolAddress string) (*models.PoolData, error) {
	if poolAddress == "" {
		return nil, ErrInvalidPoolAddress
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	poolData, err := s.dexGateway.GetPoolData(ctx, poolAddress)
	if err != nil {
		return nil, ErrFetchPoolData
	}

	return poolData, nil
}
