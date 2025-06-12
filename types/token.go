package types

import (
	"abchain_scan/chain"
	"abchain_scan/repository/orm"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"time"
)

const (
	WETH_USDC_PAIR = "0x88A43bbDF9D098eEC7bCEda4e2494615dfD9bB9C" // Uniswap v2 WETH/USDC pair
	WETH           = "0x51dA03503FBBA94B9d0D88C15690D840F02F15F4"
	USDC           = "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"
)

var (
	WETHUSDCPairAddressUniswapV2 = common.HexToAddress(WETH_USDC_PAIR)
	WETHAddress                  = common.HexToAddress(WETH)
	USDCAddress                  = common.HexToAddress(USDC)
)

func IsSameAddress(address1, address2 common.Address) bool {
	return address1.Cmp(address2) == 0
}

func IsWETH(address common.Address) bool {
	return IsSameAddress(address, WETHAddress)
}

func IsUSDC(address common.Address) bool {
	return IsSameAddress(address, USDCAddress)
}

func IsBaseToken(address common.Address) bool {
	if IsWETH(address) || IsUSDC(address) {
		return true
	}
	return false
}

type Token struct {
	Address     common.Address `json:"-"`
	Creator     common.Address `json:"-"`
	Name        string
	Symbol      string
	Decimals    int8
	TotalSupply decimal.Decimal
	BlockNumber uint64
	BlockTime   time.Time
	Program     string
	Filtered    bool
	Timestamp   time.Time
}

func (t *Token) MarshalBinary() ([]byte, error) {
	type Alias Token
	return json.Marshal(&struct {
		AddressString string `json:"Address"`
		CreatorString string `json:"Creator"`
		*Alias
	}{
		AddressString: t.Address.String(),
		CreatorString: t.Creator.String(),
		Alias:         (*Alias)(t),
	})
}

func (t *Token) UnmarshalBinary(data []byte) error {
	type Alias Token
	aux := &struct {
		AddressString string `json:"Address"`
		CreatorString string `json:"Creator"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.Address = common.HexToAddress(aux.AddressString)
	t.Creator = common.HexToAddress(aux.CreatorString)
	return nil
}

func (t *Token) Equal(token *Token) bool {
	if !IsSameAddress(t.Address, token.Address) {
		return false
	}
	if !IsSameAddress(t.Creator, token.Creator) {
		return false
	}
	if t.Name != token.Name {
		return false
	}
	if t.Symbol != token.Symbol {
		return false
	}
	if t.Decimals != token.Decimals {
		return false
	}
	if t.BlockNumber != token.BlockNumber {
		return false
	}
	if t.Program != token.Program {
		return false
	}

	return true
}

func (t *Token) GetOrmToken() *orm.Token {
	ormToken := &orm.Token{
		Address:     t.Address.String(),
		Creator:     t.Creator.String(),
		Name:        t.Name,
		Symbol:      t.Symbol,
		Decimal:     t.Decimals,
		TotalSupply: t.TotalSupply.String(),
		ChainId:     chain.Id,
		Block:       t.BlockNumber,
		BlockAt:     t.BlockTime,
		Program:     t.Program,
	}

	return ormToken.Normalize()
}
