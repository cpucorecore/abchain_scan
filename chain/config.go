package chain

import (
	"encoding/json"
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
	MainnetTerminalTotalDifficulty, _ = new(big.Int).SetString("58_750_000_000_000_000_000_000", 0)

	ChainConfig = &params.ChainConfig{
		ChainID:                 big.NewInt(Id),
		HomesteadBlock:          big.NewInt(0),
		DAOForkBlock:            big.NewInt(1_920_000),
		DAOForkSupport:          true,
		EIP150Block:             big.NewInt(2_463_000),
		EIP155Block:             big.NewInt(2_675_000),
		EIP158Block:             big.NewInt(2_675_000),
		ByzantiumBlock:          big.NewInt(4_370_000),
		ConstantinopleBlock:     big.NewInt(7_280_000),
		PetersburgBlock:         big.NewInt(7_280_000),
		IstanbulBlock:           big.NewInt(9_069_000),
		MuirGlacierBlock:        big.NewInt(9_200_000),
		BerlinBlock:             big.NewInt(12_244_000),
		LondonBlock:             big.NewInt(12_965_000),
		ArrowGlacierBlock:       big.NewInt(13_773_000),
		GrayGlacierBlock:        big.NewInt(15_050_000),
		TerminalTotalDifficulty: MainnetTerminalTotalDifficulty, // 58_750_000_000_000_000_000_000
		ShanghaiTime:            newUint64(1681338455),
		CancunTime:              newUint64(1710338135),
		PragueTime:              newUint64(1746612311),
		Ethash:                  new(params.EthashConfig),
	}
)

const (
	configJson = `
{
    "chainId": 36888,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "muirGlacierBlock": 0,
    "berlinBlock": 0,
    "clique": {
      "period": 3,
      "epoch": 30000
    }
}`
)

func init() {
	var chainConfig params.ChainConfig
	err := json.Unmarshal([]byte(configJson), &chainConfig)
	if err != nil {
		panic(err)
	}
	ChainConfig = &chainConfig
}
