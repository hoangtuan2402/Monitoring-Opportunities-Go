package gateway

import (
	"Monitoring-Opportunities/src/models"
	"context"
	"fmt"
	"math/big"
	"sync"
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

func (g *uniswapV2Gateway) GetPoolData(ctx context.Context, poolAddress string) (*models.UniswapV2Pair, error) {
	addr := common.HexToAddress(poolAddress)

	var (
		reserves    [3]*big.Int
		token0Addr  common.Address
		token1Addr  common.Address
		token0Data  *models.Token
		token1Data  *models.Token
		blockNumber uint64
		wg          sync.WaitGroup
		mu          sync.Mutex
		errs        []error
	)

	addError := func(err error) {
		if err != nil {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
		}
	}

	// Step 1: Fetch pool basic info concurrently
	wg.Add(3)
	go func() {
		defer wg.Done()
		r, err := g.getReserves(ctx, addr)
		if err == nil {
			reserves = r
		}
		addError(err)
	}()
	go func() {
		defer wg.Done()
		t0, err := g.getToken0(ctx, addr)
		if err == nil {
			token0Addr = t0
		}
		addError(err)
	}()
	go func() {
		defer wg.Done()
		t1, err := g.getToken1(ctx, addr)
		if err == nil {
			token1Addr = t1
		}
		addError(err)
	}()
	wg.Wait()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	// Step 2: Fetch token data concurrently
	wg.Add(3)
	go func() {
		defer wg.Done()
		t0, err := g.getTokenData(ctx, token0Addr)
		if err == nil {
			token0Data = t0
		}
		addError(err)
	}()
	go func() {
		defer wg.Done()
		t1, err := g.getTokenData(ctx, token1Addr)
		if err == nil {
			token1Data = t1
		}
		addError(err)
	}()
	go func() {
		defer wg.Done()
		bn, err := g.ethClient.BlockNumber(ctx)
		if err == nil {
			blockNumber = bn
		}
		addError(err)
	}()
	wg.Wait()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	return &models.UniswapV2Pair{
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
	var (
		symbol   string
		decimals int
		wg       sync.WaitGroup
		mu       sync.Mutex
		errs     []error
	)

	addError := func(err error) {
		if err != nil {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
		}
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		s, err := g.getTokenSymbol(ctx, tokenAddr)
		if err == nil {
			symbol = s
		}
		addError(err)
	}()
	go func() {
		defer wg.Done()
		d, err := g.getTokenDecimals(ctx, tokenAddr)
		if err == nil {
			decimals = d
		}
		addError(err)
	}()
	wg.Wait()

	if len(errs) > 0 {
		return nil, errs[0]
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
