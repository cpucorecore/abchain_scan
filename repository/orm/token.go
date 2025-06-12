package orm

import (
	"abchain_scan/util"
	"time"
	"unicode/utf8"
)

type Token struct {
	Address     string
	Creator     string
	Name        string
	Symbol      string
	Decimal     int8
	TotalSupply string
	ChainId     int
	Block       uint64
	BlockAt     time.Time
	Program     string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	MainPair    string
}

func (t *Token) TableName() string {
	return "token"
}

func (t *Token) Equal(t2 *Token) bool {
	if t.Address != t2.Address {
		return false
	}
	if t.Name != t2.Name {
		return false
	}
	if t.Symbol != t2.Symbol {
		return false
	}
	if t.Decimal != t2.Decimal {
		return false
	}
	if t.TotalSupply != t2.TotalSupply {
		return false
	}
	if t.Block != t2.Block {
		return false
	}
	return true
}

func (t *Token) Normalize() *Token {
	const (
		maxNameLength   = 64 // Maximum allowed characters for name
		maxSymbolLength = 32 // Maximum allowed characters for symbol
		maxSupplyLength = 64 // Maximum allowed characters for total supply
	)

	if utf8.RuneCountInString(t.Name) > maxNameLength {
		t.Name = util.TruncateToMaxChars(t.Name, maxNameLength)
	}

	if utf8.RuneCountInString(t.Symbol) > maxSymbolLength {
		t.Symbol = util.TruncateToMaxChars(t.Symbol, maxSymbolLength)
	}

	if utf8.RuneCountInString(t.TotalSupply) > maxSupplyLength {
		t.TotalSupply = "0"
	}
	return t
}
