package chain

import (
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

/*
from geth v1.15.11
TODO: should use op-geth chain config:
https://github.com/ethereum-optimism/op-geth/blob/optimism/params/config.go
*/

func newUint64(v uint64) *uint64 {
	return &v
}

var (
	ChainConfig = &params.ChainConfig{
		ChainID:             big.NewInt(Id),
		HomesteadBlock:      big.NewInt(1_150_000),
		DAOForkBlock:        big.NewInt(1_920_000),
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(2_463_000),
		EIP155Block:         big.NewInt(2_675_000),
		EIP158Block:         big.NewInt(2_675_000),
		ByzantiumBlock:      big.NewInt(4_370_000),
		ConstantinopleBlock: big.NewInt(7_280_000),
		PetersburgBlock:     big.NewInt(7_280_000),
		IstanbulBlock:       big.NewInt(9_069_000),
		MuirGlacierBlock:    big.NewInt(9_200_000),
		Ethash:              new(params.EthashConfig),
	}
)
