package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestToken_MarshalBinary(t *testing.T) {
	address := common.HexToAddress("0xE76004cFFcAb665C4692F663B8FB2A2F66AdDa9B")
	token := Token{
		Address:     address,
		Creator:     address,
		Name:        "test",
		Symbol:      "test",
		Decimals:    18,
		TotalSupply: decimal.NewFromInt(1),
		BlockNumber: 1,
		BlockTime:   time.Unix(1000, 1),
		Program:     "test",
	}

	tokenBytes, err := token.MarshalBinary()
	require.NoError(t, err)
	t.Log(string(tokenBytes))

	token2 := &Token{}
	err = token2.UnmarshalBinary(tokenBytes)
	require.NoError(t, err)
	require.True(t, token.Equal(token2))
}
