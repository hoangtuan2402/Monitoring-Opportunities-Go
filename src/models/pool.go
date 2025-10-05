package models

import "time"

type Token struct {
	Address  string `json:"address"`
	Decimals int    `json:"decimals"`
	Symbol   string `json:"symbol"`
}

type UniswapV2Pair struct {
	Address            string    `json:"address"`
	Reserve0           string    `json:"reserve0"`
	Reserve1           string    `json:"reserve1"`
	Token0             Token     `json:"token0"`
	Token1             Token     `json:"token1"`
	BlockTimestampLast int64     `json:"blockTimestampLast"`
	DownloadedAt       time.Time `json:"downloadedAt"`
	Network            string    `json:"network"`
	BlockNumber        int64     `json:"blockNumber"`
	Symbol             string    `json:"symbol"`
	Type               string    `json:"type"`
}
