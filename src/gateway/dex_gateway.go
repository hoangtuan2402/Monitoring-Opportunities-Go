package gateway

import (
	"Monitoring-Opportunities/src/models"
	"context"
)

type DEXGateway interface {
	GetPoolData(ctx context.Context, poolAddress string) (*models.PoolData, error)
}
