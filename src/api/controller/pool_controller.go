package handler

import (
	"Monitoring-Opportunities/src/common"
	"Monitoring-Opportunities/src/models"
	service "Monitoring-Opportunities/src/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PoolController struct {
	poolService service.PoolService
}

func NewPoolController(poolService service.PoolService) *PoolController {
	return &PoolController{
		poolService: poolService,
	}
}

func (c *PoolController) GetPoolData(ctx *gin.Context) {
	poolAddress := ctx.Query("address")
	if poolAddress == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "pool address is required"})
		return
	}

	poolData, err := c.poolService.GetPoolData(poolAddress)
	if err != nil {
		log.Printf("Failed to get pool data for %s: %v", poolAddress, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	ctx.JSON(
		http.StatusOK,
		common.BaseResponse[*models.PoolData]{
			Status:  http.StatusOK,
			Message: "Successfully retrieved pool data",
			Data:    poolData,
		},
	)
}
