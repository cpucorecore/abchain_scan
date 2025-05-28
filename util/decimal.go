package util

import "github.com/shopspring/decimal"

var (
	Epsilon = decimal.NewFromFloat(1e-12)
)

func DecimalEqual(a, b decimal.Decimal) bool {
	return a.Sub(b).Abs().LessThanOrEqual(Epsilon)
}
