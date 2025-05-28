package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPair_IsFiltered(t *testing.T) {
	pair := &Pair{
		Address:    common.HexToAddress("0xF6C8490Df6a5bFCc07484DC87254B4139C9CCCd3"),
		Token0Core: &TokenCore{Address: common.HexToAddress("0xD1E0f3957E91282Bc1acB95fFaaDBa58ac11BeeD")},
		Token1Core: &TokenCore{Address: common.HexToAddress("0xfe0A4739139D5b64b9fA86DA767B464086A9d5B2")},
	}

	filtered := pair.FilterByToken0AndToken1()
	assert.Equal(t, filtered, true)
	assert.Equal(t, pair.Filtered, true)
	assert.Equal(t, pair.FilterCode, FilterCodeNoBaseToken)
}

func TestPair_OrderTokens(t *testing.T) {
	TokenWETH := &TokenCore{
		Address: WETHAddress,
	}
	TokenUSDC := &TokenCore{
		Address: USDCAddress,
	}
	TokenNonBase := &TokenCore{
		Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
	}

	tests := []struct {
		name           string
		pair           *Pair
		expected       *Pair
		TokensReversed bool
	}{
		{
			name: "WETH/USDC",
			pair: &Pair{
				Token0Core: TokenWETH,
				Token1Core: TokenUSDC,
			},
			expected: &Pair{
				Token0Core: TokenWETH,
				Token1Core: TokenUSDC,
			},
			TokensReversed: false,
		},
		{
			name: "USDC/WETH",
			pair: &Pair{
				Token0Core: TokenUSDC,
				Token1Core: TokenWETH,
			},
			expected: &Pair{
				Token0Core: TokenWETH,
				Token1Core: TokenUSDC,
			},
			TokensReversed: true,
		},
		{
			name: "WETH/NonBase",
			pair: &Pair{
				Token0Core: TokenWETH,
				Token1Core: TokenNonBase,
			},
			expected: &Pair{
				Token0Core: TokenNonBase,
				Token1Core: TokenWETH,
			},
			TokensReversed: true,
		},
		{
			name: "NonBase/WETH",
			pair: &Pair{
				Token0Core: TokenNonBase,
				Token1Core: TokenWETH,
			},
			expected: &Pair{
				Token0Core: TokenNonBase,
				Token1Core: TokenWETH,
			},
			TokensReversed: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.pair.OrderToken0Token1()
			assert.Equal(t, test.pair.Token0Core.Address, test.expected.Token0Core.Address)
			assert.Equal(t, test.pair.Token1Core.Address, test.expected.Token1Core.Address)
			assert.Equal(t, test.pair.TokensReversed, test.TokensReversed)
		})
	}
}

func TestPair_MarshalBinary(t *testing.T) {
	address := common.HexToAddress("0xE76004cFFcAb665C4692F663B8FB2A2F66AdDa9B")
	token := &Token{
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

	tokenBasic := &TokenCore{
		Address:  address,
		Symbol:   "test",
		Decimals: 18,
	}

	pair := &Pair{
		Address:          address,
		TokensReversed:   false,
		Token0Core:       tokenBasic,
		Token1Core:       tokenBasic,
		Token0:           token,
		Token1:           token,
		Token0InitAmount: decimal.NewFromInt(1),
		Token1InitAmount: decimal.NewFromInt(1),
		Block:            1,
		ProtocolId:       1,
		Filtered:         true,
		FilterCode:       1,
	}

	bytes, err := pair.MarshalBinary()
	require.NoError(t, err)
	t.Log(string(bytes))

	pair2 := &Pair{}
	err = pair2.UnmarshalBinary(bytes)
	require.NoError(t, err)
	require.True(t, pair.Equal(pair2))
}
