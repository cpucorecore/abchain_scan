package block_getter

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func Test_GetBlock(t *testing.T) {
	ethCli, err := ethclient.Dial("http://47.251.86.106:8808")
	require.NoError(t, err)
	block, getBlockErr := ethCli.BlockByNumber(context.Background(), big.NewInt(67488764))
	require.NoError(t, getBlockErr)
	t.Log(block)
}
