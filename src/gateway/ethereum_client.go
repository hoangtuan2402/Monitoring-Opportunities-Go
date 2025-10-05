package gateway

import (
	"Monitoring-Opportunities/src/config"
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	CallContract(ctx context.Context, to common.Address, data []byte) ([]byte, error)
}

type ethereumClient struct {
	client *ethclient.Client
}

func NewEthereumRpcClient(cfg config.Config) (EthereumClient, error) {
	client, err := ethclient.Dial(cfg.EthRPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	return &ethereumClient{client: client}, nil
}

func (c *ethereumClient) BlockNumber(ctx context.Context) (uint64, error) {
	return c.client.BlockNumber(ctx)
}

func (c *ethereumClient) CallContract(ctx context.Context, to common.Address, data []byte) ([]byte, error) {
	callMsg := ethereum.CallMsg{
		To:   &to,
		Data: data,
	}

	return c.client.CallContract(ctx, callMsg, nil)
}

func EncodeMethodCall(signature string) []byte {
	return crypto.Keccak256([]byte(signature))[:4]
}

func DecodeString(data []byte) (string, error) {
	stringType, _ := abi.NewType("string", "", nil)
	decoded, err := abi.Arguments{{Type: stringType}}.Unpack(data)
	if err != nil {
		return "", err
	}
	return decoded[0].(string), nil
}

func DecodeBigInt(data []byte) *big.Int {
	return new(big.Int).SetBytes(data)
}

func DecodeAddress(data []byte) common.Address {
	return common.BytesToAddress(data)
}
