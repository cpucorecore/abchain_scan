package event

import (
	"abchain_scan/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"math/big"
)

func ParseAmountsByPair(token0AmountWei, token1AmountWei *big.Int, pair *types.Pair) (token0Amount, token1Amount decimal.Decimal) {
	if !pair.TokensReversed {
		token0Amount = decimal.NewFromBigInt(token0AmountWei, -(int32)(pair.Token0Core.Decimals))
		token1Amount = decimal.NewFromBigInt(token1AmountWei, -(int32)(pair.Token1Core.Decimals))
	} else {
		token0Amount = decimal.NewFromBigInt(token1AmountWei, -(int32)(pair.Token0Core.Decimals))
		token1Amount = decimal.NewFromBigInt(token0AmountWei, -(int32)(pair.Token1Core.Decimals))
	}
	return
}

func CalcAmountAndPrice(
	bnbPrice decimal.Decimal,
	token0Amount, token1Amount decimal.Decimal,
	token1Address common.Address,
) (amountUSD, priceUSD decimal.Decimal) {
	if types.IsWETH(token1Address) {
		amountUSD = token1Amount.Mul(bnbPrice)
		if !token0Amount.IsZero() {
			priceUSD = amountUSD.Div(token0Amount)
		}
	} else if types.IsUSDC(token1Address) {
		amountUSD = token1Amount
		if !token0Amount.IsZero() {
			priceUSD = amountUSD.Div(token0Amount)
		}
	}
	return
}
