package gateway

import (
	"Monitoring-Opportunities/src/models"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type uniswapV2Gateway struct {
	ethClient EthereumClient
}

func NewUniswapV2Gateway(ethClient EthereumClient) DEXGateway {
	return &uniswapV2Gateway{
		ethClient: ethClient,
	}
}

func (g *uniswapV2Gateway) GetPoolData(ctx context.Context, poolAddress string) (*models.PoolData, error) {
	addr := common.HexToAddress(poolAddress)

	reserves, err := g.getReserves(ctx, addr)
	if err != nil {
		return nil, err
	}

	token0Addr, err := g.getToken0(ctx, addr)
	if err != nil {
		return nil, err
	}

	token1Addr, err := g.getToken1(ctx, addr)
	if err != nil {
		return nil, err
	}

	token0Data, err := g.getTokenData(ctx, token0Addr)
	if err != nil {
		return nil, err
	}

	token1Data, err := g.getTokenData(ctx, token1Addr)
	if err != nil {
		return nil, err
	}

	blockNumber, err := g.ethClient.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get block number: %w", err)
	}

	return &models.PoolData{
		Address:            poolAddress,
		Reserve0:           reserves[0].String(),
		Reserve1:           reserves[1].String(),
		Token0:             *token0Data,
		Token1:             *token1Data,
		BlockTimestampLast: reserves[2].Int64(),
		DownloadedAt:       time.Now(),
		Network:            "mainnet",
		BlockNumber:        int64(blockNumber),
		Symbol:             fmt.Sprintf("%s/%s", token0Data.Symbol, token1Data.Symbol),
		Type:               "v2",
	}, nil
}

func (g *uniswapV2Gateway) getReserves(ctx context.Context, poolAddr common.Address) ([3]*big.Int, error) {
	data := EncodeMethodCall("getReserves()")

	result, err := g.ethClient.CallContract(ctx, poolAddr, data)
	if err != nil {
		return [3]*big.Int{}, fmt.Errorf("failed to call getReserves: %w", err)
	}

	reserve0 := DecodeBigInt(result[0:32])
	reserve1 := DecodeBigInt(result[32:64])
	blockTimestamp := DecodeBigInt(result[64:96])

	return [3]*big.Int{reserve0, reserve1, blockTimestamp}, nil
}

func (g *uniswapV2Gateway) getToken0(ctx context.Context, poolAddr common.Address) (common.Address, error) {
	data := EncodeMethodCall("token0()")

	result, err := g.ethClient.CallContract(ctx, poolAddr, data)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to call token0: %w", err)
	}

	return DecodeAddress(result), nil
}

func (g *uniswapV2Gateway) getToken1(ctx context.Context, poolAddr common.Address) (common.Address, error) {
	data := EncodeMethodCall("token1()")

	result, err := g.ethClient.CallContract(ctx, poolAddr, data)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to call token1: %w", err)
	}

	return DecodeAddress(result), nil
}

func (g *uniswapV2Gateway) getTokenData(ctx context.Context, tokenAddr common.Address) (*models.Token, error) {
	symbol, err := g.getTokenSymbol(ctx, tokenAddr)
	if err != nil {
		return nil, err
	}

	decimals, err := g.getTokenDecimals(ctx, tokenAddr)
	if err != nil {
		return nil, err
	}

	return &models.Token{
		Address:  tokenAddr.Hex(),
		Symbol:   symbol,
		Decimals: decimals,
	}, nil
}

func (g *uniswapV2Gateway) getTokenSymbol(ctx context.Context, tokenAddr common.Address) (string, error) {
	data := EncodeMethodCall("symbol()")

	result, err := g.ethClient.CallContract(ctx, tokenAddr, data)
	if err != nil {
		return "", fmt.Errorf("failed to call symbol: %w", err)
	}

	return DecodeString(result)
}

func (g *uniswapV2Gateway) getTokenDecimals(ctx context.Context, tokenAddr common.Address) (int, error) {
	data := EncodeMethodCall("decimals()")

	result, err := g.ethClient.CallContract(ctx, tokenAddr, data)
	if err != nil {
		return 0, fmt.Errorf("failed to call decimals: %w", err)
	}

	decimals := DecodeBigInt(result)
	return int(decimals.Int64()), nil
}
