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

type UniswapV2PairService interface {
	GetPoolData(poolAddress string) (*models.UniswapV2Pair, error)
}

type uniswapV2PairService struct {
	dexGateway gateway.DEXGateway
}

func NewUniswapV2PairService(dexGateway gateway.DEXGateway) UniswapV2PairService {
	return &uniswapV2PairService{
		dexGateway: dexGateway,
	}
}

func (s *uniswapV2PairService) GetPoolData(poolAddress string) (*models.UniswapV2Pair, error) {
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
